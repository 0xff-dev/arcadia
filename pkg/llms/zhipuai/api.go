/*
Copyright 2023 KubeAGI.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// NOTE: Reference zhipuai's python sdk: model_api/api.py

package zhipuai

import (
	"errors"
	"fmt"
	"time"

	"github.com/kubeagi/arcadia/pkg/llms"
	"github.com/r3labs/sse/v2"
)

const (
	ZhipuaiModelAPIURL         = "https://open.bigmodel.cn/api/paas/v3/model-api"
	ZhipuaiModelDefaultTimeout = 300 * time.Second
)

type Model string

const (
	ZhiPuAILite Model = "chatglm_lite"
	ZhiPuAIStd  Model = "chatglm_std"
	ZhiPuAIPro  Model = "chatglm_pro"
)

type Method string

const (
	// POST
	ZhiPuAIInvoke      Method = "invoke"
	ZhiPuAIAsyncInvoke Method = "async-invoke"
	ZhiPuAISSEInvoke   Method = "sse-invoke"
	// GET
	ZhiPuAIAsyncGet Method = "async-get"
)

func BuildAPIURL(model Model, method Method) string {
	return fmt.Sprintf("%s/%s/%s", ZhipuaiModelAPIURL, model, method)
}

var _ llms.LLM = (*ZhiPuAI)(nil)

type ZhiPuAI struct {
	apiKey string
}

func NewZhiPuAI(apiKey string) *ZhiPuAI {
	return &ZhiPuAI{
		apiKey: apiKey,
	}
}

func (z ZhiPuAI) Type() llms.LLMType {
	return llms.ZhiPuAI
}

// Call wraps a common AI api call
func (z *ZhiPuAI) Call(params ModelParams) (*Response, error) {
	switch params.Method {
	case ZhiPuAIInvoke:
		return z.Invoke(params)
	case ZhiPuAIAsyncInvoke:
		return z.AsyncInvoke(params)
	case ZhiPuAIAsyncGet:
		return z.Get(params)
	default:
		return nil, errors.New("unknown method")
	}
}

// Invoke calls zhipuai and returns result immediately
func (z *ZhiPuAI) Invoke(params ModelParams) (*Response, error) {
	url := BuildAPIURL(params.Model, ZhiPuAIInvoke)
	token, err := GenerateToken(z.apiKey, APITokenTTLSeconds)
	if err != nil {
		return nil, err
	}

	return Post(url, token, params, ZhipuaiModelDefaultTimeout)
}

// AsyncInvoke only returns a task id which can be used to get result of task later
func (z *ZhiPuAI) AsyncInvoke(params ModelParams) (*Response, error) {
	url := BuildAPIURL(params.Model, ZhiPuAIAsyncInvoke)
	token, err := GenerateToken(z.apiKey, APITokenTTLSeconds)
	if err != nil {
		return nil, err
	}

	return Post(url, token, params, ZhipuaiModelDefaultTimeout)
}

// Get result of task async-invoke
func (z *ZhiPuAI) Get(params ModelParams) (*Response, error) {
	if params.TaskID == "" {
		return nil, errors.New("TaskID is required when running Get with method AsyncInvoke")
	}

	// url with task id
	url := fmt.Sprintf("%s/%s", BuildAPIURL(params.Model, ZhiPuAIAsyncInvoke), params.TaskID)
	token, err := GenerateToken(z.apiKey, APITokenTTLSeconds)
	if err != nil {
		return nil, err
	}

	return Get(url, token, ZhipuaiModelDefaultTimeout)
}

func (z *ZhiPuAI) SSEInvoke(params ModelParams, handler func(*sse.Event)) error {
	url := BuildAPIURL(params.Model, ZhiPuAISSEInvoke)
	token, err := GenerateToken(z.apiKey, APITokenTTLSeconds)
	if err != nil {
		return err
	}
	return Stream(url, token, params, ZhipuaiModelDefaultTimeout, nil)
}

func (z *ZhiPuAI) Validate() (llms.Response, error) {
	url := BuildAPIURL(ZhiPuAILite, ZhiPuAIInvoke)
	token, err := GenerateToken(z.apiKey, APITokenTTLSeconds)
	if err != nil {
		return nil, err
	}

	testPrompt := []Prompt{
		{
			Role:    "user",
			Content: "Hello!",
		},
	}

	testParam := ModelParams{
		Method:      ZhiPuAIAsyncInvoke,
		Model:       ZhiPuAILite,
		Temperature: 0.95,
		TopP:        0.7,
		Prompt:      testPrompt,
	}

	postResponse, err := Post(url, token, testParam, ZhipuaiModelDefaultTimeout)
	if err != nil {
		return nil, err
	}

	return postResponse, nil
}