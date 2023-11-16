/*
 (c) Copyright [2023] Open Text.
 Licensed under the Apache License, Version 2.0 (the "License");
 You may not use this file except in compliance with the License.
 You may obtain a copy of the License at

 http://www.apache.org/licenses/LICENSE-2.0

 Unless required by applicable law or agreed to in writing, software
 distributed under the License is distributed on an "AS IS" BASIS,
 WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 See the License for the specific language governing permissions and
 limitations under the License.
*/

package vclusterops

import (
	"encoding/json"
	"fmt"

	"github.com/vertica/vcluster/vclusterops/vlog"
)

type NMAStageErrorReportOp struct {
	ScrutinizeOpBase
}

type stageErrorReportRequestData struct {
	CatalogPath string `json:"catalog_path"`
}

type stageErrorReportResponseData struct {
	Name      string `json:"name"`
	SizeBytes int64  `json:"size_bytes"`
	ModTime   string `json:"mod_time"`
}

func makeNMAStageErrorReportOp(logger vlog.Printer,
	id string,
	hosts []string,
	hostNodeNameMap map[string]string,
	hostCatPathMap map[string]string) (NMAStageErrorReportOp, error) {
	// base members
	op := NMAStageErrorReportOp{}
	op.name = "NMAStageErrorReportOp"
	op.logger = logger.WithName(op.name)
	op.hosts = hosts

	// scrutinize members
	op.id = id
	op.batch = scrutinizeBatchContext
	op.hostNodeNameMap = hostNodeNameMap
	op.hostCatPathMap = hostCatPathMap
	op.httpMethod = PostMethod
	op.urlSuffix = "/ErrorReport.txt"

	// the caller is responsible for making sure hosts and maps match up exactly
	err := validateHostMaps(hosts, hostNodeNameMap, hostCatPathMap)
	return op, err
}

func (op *NMAStageErrorReportOp) setupRequestBody(hosts []string) error {
	op.hostRequestBodyMap = make(map[string]string, len(hosts))
	for _, host := range hosts {
		stageErrorReportData := stageErrorReportRequestData{}
		stageErrorReportData.CatalogPath = op.hostCatPathMap[host]

		dataBytes, err := json.Marshal(stageErrorReportData)
		if err != nil {
			return fmt.Errorf("[%s] fail to marshal request data to JSON string, detail %w", op.name, err)
		}

		op.hostRequestBodyMap[host] = string(dataBytes)
	}

	return nil
}

func (op *NMAStageErrorReportOp) prepare(execContext *OpEngineExecContext) error {
	err := op.setupRequestBody(op.hosts)
	if err != nil {
		return err
	}
	execContext.dispatcher.setup(op.hosts)

	return op.setupClusterHTTPRequest(op.hosts)
}

func (op *NMAStageErrorReportOp) execute(execContext *OpEngineExecContext) error {
	if err := op.runExecute(execContext); err != nil {
		return err
	}

	return op.processResult(execContext)
}

func (op *NMAStageErrorReportOp) finalize(_ *OpEngineExecContext) error {
	return nil
}

func (op *NMAStageErrorReportOp) processResult(_ *OpEngineExecContext) error {
	fileList := make([]stageErrorReportResponseData, 0)
	return processStagedFilesResult(&op.ScrutinizeOpBase, fileList)
}