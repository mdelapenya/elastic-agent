// Licensed to Elasticsearch B.V. under one or more contributor
// license agreements. See the NOTICE file distributed with
// this work for additional information regarding copyright
// ownership. Elasticsearch B.V. licenses this file to you under
// the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.

package handlers

import (
	"context"
	"fmt"
	"time"

	"github.com/elastic/elastic-agent-poc/elastic-agent/pkg/agent/errors"
	"github.com/elastic/elastic-agent-poc/elastic-agent/pkg/agent/storage/store"
	"github.com/elastic/elastic-agent-poc/elastic-agent/pkg/core/logger"
	"github.com/elastic/elastic-agent-poc/elastic-agent/pkg/core/server"
	"github.com/elastic/elastic-agent-poc/elastic-agent/pkg/fleetapi"
)

const (
	defaultActionTimeout = time.Minute
	maxActionTimeout     = time.Hour
)

var errActionTimeoutInvalid = errors.New("action timeout is invalid")

// AppAction is a handler for application actions.
type AppAction struct {
	log *logger.Logger
	srv *server.Server
}

// NewAppAction creates a new AppAction handler.
func NewAppAction(log *logger.Logger, srv *server.Server) *AppAction {
	return &AppAction{
		log: log,
		srv: srv,
	}
}

// Handle handles application action.
func (h *AppAction) Handle(ctx context.Context, a fleetapi.Action, acker store.FleetAcker) error {
	h.log.Debugf("handlerAppAction: action '%+v' received", a)
	action, ok := a.(*fleetapi.ActionApp)
	if !ok {
		return fmt.Errorf("invalid type, expected ActionApp and received %T", a)
	}

	appState, ok := h.srv.FindByInputType(action.InputType)
	if !ok {
		// If the matching action is not found ack the action with the error for action result document
		action.StartedAt = time.Now().UTC().Format(time.RFC3339Nano)
		action.CompletedAt = action.StartedAt
		action.Error = fmt.Sprintf("matching app is not found for action input: %s", action.InputType)
		return acker.Ack(ctx, action)
	}

	params, err := action.MarshalMap()
	if err != nil {
		return err
	}

	start := time.Now().UTC()
	timeout := defaultActionTimeout
	if action.Timeout > 0 {
		timeout = time.Duration(action.Timeout) * time.Second
		if timeout > maxActionTimeout {
			h.log.Debugf("handlerAppAction: action '%v' timeout exceeds maximum allowed %v", action.InputType, maxActionTimeout)
			err = errActionTimeoutInvalid
		}
	}

	var res map[string]interface{}
	if err == nil {
		h.log.Debugf("handlerAppAction: action '%v' started with timeout: %v", action.InputType, timeout)
		res, err = appState.PerformAction(action.InputType, params, timeout)
	}
	end := time.Now().UTC()

	startFormatted := start.Format(time.RFC3339Nano)
	endFormatted := end.Format(time.RFC3339Nano)
	h.log.Debugf("handlerAppAction: action '%v' finished, startFormatted: %v, endFormatted: %v, err: %v", action.InputType, startFormatted, endFormatted, err)
	if err != nil {
		action.StartedAt = startFormatted
		action.CompletedAt = endFormatted
		action.Error = err.Error()
	} else {
		action.StartedAt = readMapString(res, "started_at", startFormatted)
		action.CompletedAt = readMapString(res, "completed_at", endFormatted)
		action.Error = readMapString(res, "error", "")
		appendActionResponse(action, action.InputType, res)
	}

	return acker.Ack(ctx, action)
}

var (
	none = struct{}{}

	// The set of action response fields are not included in the action_response property, because there are already set to top level fields
	excludeActionResponseFields = map[string]struct{}{
		"started_at":   none,
		"completed_at": none,
		"error":        none,
	}
)

// appendActionResponse appends the action response property with all the action response values excluding the ones specified in excludeActionResponseFields
// "action_response": {
// 	   "endpoint": {
// 		   "acked": true
// 	   }
//  }
func appendActionResponse(action *fleetapi.ActionApp, inputType string, res map[string]interface{}) {
	if len(res) == 0 {
		return
	}

	m := make(map[string]interface{}, len(res))

	for k, v := range res {
		if _, ok := excludeActionResponseFields[k]; !ok {
			m[k] = v
		}
	}

	if len(m) > 0 {
		mt := make(map[string]interface{}, 1)
		mt[inputType] = m

		action.Response = mt
	}
}

func readMapString(m map[string]interface{}, key string, def string) string {
	if m == nil {
		return def
	}

	if v, ok := m[key]; ok {
		if s, ok := v.(string); ok && s != "" {
			return s
		}
	}
	return def
}
