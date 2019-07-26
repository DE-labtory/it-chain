/*
 * Copyright 2018 DE-labtory
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 * https://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package tesseract

import (
	"encoding/json"
	"errors"
	"fmt"
	"sync"

	"path"

	"net"
	"strconv"

	"strings"

	"github.com/DE-labtory/iLogger"
	"github.com/DE-labtory/it-chain/ivm"
	"github.com/DE-labtory/tesseract"
	"github.com/DE-labtory/tesseract/container"
	"github.com/DE-labtory/tesseract/pb"
	"github.com/rs/xid"
)

var ErrContainerDoesNotExist = errors.New("container does not exist")
var ErrICodeInfoMapNotEmpty = errors.New("ICode info struct in current container is not empty")

type ICodeInfo struct {
	container     tesseract.Container
	containerIp   string
	containerPort string
	iCode         ivm.ICode
}

type ContainerService struct {
	sync.RWMutex
	iCodeInfoMap map[tesseract.ContainerID]ICodeInfo
	dockerConfig *ContainerDockerConfig
}
type ContainerDockerConfig struct {
	Subnet      string
	VolumeName  string
	NetworkName string
}

func NewContainerService(config *ContainerDockerConfig) (*ContainerService, error) {
	return &ContainerService{
		dockerConfig: config,
		iCodeInfoMap: make(map[tesseract.ContainerID]ICodeInfo),
		RWMutex:      sync.RWMutex{},
	}, nil
}

func (cs ContainerService) StartContainer(icode ivm.ICode) error {
	iLogger.Info(nil, fmt.Sprintf("[IVM] Starting container - icodeID: [%s]", icode.ID))

	cs.Lock()
	defer cs.Unlock()

	conf, err := cs.createContainerConfig(icode)
	if err != nil {
		return err
	}

	createdContainer, err := container.Create(conf)
	if err != nil {
		return err
	}

	_, ok := cs.iCodeInfoMap[icode.ID]
	if !ok {
		iCodeInfo := ICodeInfo{
			container: createdContainer,
			iCode:     icode,
		}

		cs.iCodeInfoMap[icode.ID] = iCodeInfo
	} else {
		return ErrICodeInfoMapNotEmpty
	}

	return nil
}

func (cs ContainerService) createContainerConfig(icode ivm.ICode) (tesseract.ContainerConfig, error) {
	switch cs.hasDockerConfig() {
	case true:
		return cs.createDockerConfig(icode)
	case false:
		return cs.createHostMachineConfig(icode)
	}

	return tesseract.ContainerConfig{}, nil
}

func (cs ContainerService) hasDockerConfig() bool {
	return cs.dockerConfig != nil
}

func (cs ContainerService) createDockerConfig(icode ivm.ICode) (tesseract.ContainerConfig, error) {
	subnetRootIp, err := getSubnetRootIp(cs.dockerConfig.Subnet)
	if err != nil {
		return tesseract.ContainerConfig{}, err
	}

	hostIp := subnetRootIp + strconv.Itoa(2)

	port, err := getPort(hostIp, "60000")
	if err != nil {
		return tesseract.ContainerConfig{}, err
	}

	containerIp := subnetRootIp + strconv.Itoa(len(cs.iCodeInfoMap)+3)

	return tesseract.ContainerConfig{
		Language: "go",
		Name:     icode.ID,
		ContainerImage: tesseract.ContainerImage{
			Name: "golang",
			Tag:  "1.9",
		},
		HostIp:      hostIp,
		Port:        port,
		ContainerIp: containerIp,
		StartCmd:    []string{"go", "run", path.Join("/go/src", icode.FolderName, "icode.go"), "-p", port},
		Mount:       []string{cs.dockerConfig.VolumeName + ":" + "/go/src"},
		Network: &tesseract.Network{
			Name: cs.dockerConfig.NetworkName,
		},
	}, nil
}

func (cs ContainerService) createHostMachineConfig(icode ivm.ICode) (tesseract.ContainerConfig, error) {
	hostIp := "0.0.0.0"

	port, err := container.GetAvailablePort()
	if err != nil {
		return tesseract.ContainerConfig{}, err
	}

	return tesseract.ContainerConfig{
		Language: "go",
		Name:     icode.ID,
		ContainerImage: tesseract.ContainerImage{
			Name: "golang",
			Tag:  "1.9",
		},
		HostIp:   hostIp,
		Port:     port,
		StartCmd: []string{"go", "run", path.Join("/go/src", icode.FolderName, "icode.go"), "-p", port},
		Mount:    []string{icode.Path + ":" + "/go/src/" + icode.FolderName},
	}, nil
}

func getPort(findIp string, startPort string) (string, error) {
	tryToFindPort, err := strconv.Atoi(startPort)
	if err != nil {
		return "", err
	}
	for {
		if st, _ := strconv.Atoi(startPort); tryToFindPort-st > 10000 { // 최대 포트검색시도는 10000개
			return "", errors.New("cant find empty error while 10000")
		}
		lis, err := net.Listen("tcp", findIp+":"+strconv.Itoa(tryToFindPort))
		lis.Close()
		if err != nil {
			continue
		}
		return strconv.Itoa(tryToFindPort), nil
	}
}

func (cs ContainerService) ExecuteRequest(request ivm.Request) (ivm.Result, error) {
	iLogger.Info(nil, fmt.Sprintf("[IVM] Executing icode - icodeID: [%s]", request.ICodeID))

	iCodeInfo, ok := cs.iCodeInfoMap[request.ICodeID]
	if !ok {
		return ivm.Result{}, ErrContainerDoesNotExist
	}

	resultCh := make(chan ivm.Result)
	errCh := make(chan error)

	var callback = func(response *pb.Response, err error) {
		if err != nil {
			errCh <- err
		}

		data := make(map[string]string)
		if len(response.Data) != 0 {
			if err = json.Unmarshal(response.Data, &data); err != nil {
				errCh <- err
				return
			}
		}

		resultCh <- ivm.Result{
			Err:  response.Error,
			Data: data,
		}
	}

	err := iCodeInfo.container.Request(tesseract.Request{
		Uuid:     xid.New().String(),
		Args:     request.Args,
		FuncName: request.Function,
		TypeName: request.Type,
	}, callback)

	if err != nil {
		errCh <- err
	}

	select {
	case err := <-errCh:
		iLogger.Error(nil, fmt.Sprintf("[IVM] fail executing ivm, id:%s", request.ICodeID))
		return ivm.Result{}, err
	case result := <-resultCh:
		return result, nil
	}
}

func (cs ContainerService) StopContainer(id ivm.ID) error {
	iCodeInfo, ok := cs.iCodeInfoMap[id]

	if !ok {
		return ErrContainerDoesNotExist
	}

	err := iCodeInfo.container.Close()

	if err != nil {
		return err
	}

	delete(cs.iCodeInfoMap, id)
	return nil
}

func (cs ContainerService) GetRunningICodeList() []ivm.ICode {
	iCodeList := make([]ivm.ICode, 0)

	for id := range cs.iCodeInfoMap {
		iCodeInfo, ok := cs.iCodeInfoMap[id]
		if ok {
			iCodeList = append(iCodeList, iCodeInfo.iCode)
		}
	}

	return iCodeList
}

func getSubnetRootIp(subnet string) (string, error) {
	if !strings.Contains(subnet, "/") {
		return "", errors.New("is not formatted for subnet")
	}
	args := strings.Split(subnet, "/")
	if len(args) != 2 {
		return "", errors.New("is not formatted for subnet")
	}
	//172.88.x.0/8
	subnetCount, err := strconv.Atoi(args[1])
	if err != nil {
		return "", err
	}
	if subnetCount%8 != 0 {
		return "", errors.New("subnet 비트는 8배수여야함")
	}
	subnetCount /= 8
	dotSplittedIp := strings.Split(args[0], ".")
	rootIp := ""
	for i := 0; i < subnetCount; i++ {
		rootIp = rootIp + dotSplittedIp[i] + "."
	}
	return rootIp, nil
}
