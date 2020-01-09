/*
 * Copyright 2020 Nalej
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package utils

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/nalej/system-model/internal/pkg/entities"
	"os"
)

func RunIntegrationTests() bool {
	var runIntegration = os.Getenv("RUN_INTEGRATION_TEST")
	return runIntegration == "true"
}


func CreateOrganization() *entities.Organization {
	return entities.NewOrganization(fmt.Sprintf("Nalej-%s", uuid.New().String()), "Nalej Test Address", "City Test", "State Test", "U.S.A", "XXX", "Photo")
}