/*
 *    Copyright 2022 chenquan
 *
 *    Licensed under the Apache License, Version 2.0 (the "License");
 *    you may not use this file except in compliance with the License.
 *    You may obtain a copy of the License at
 *
 *        http://www.apache.org/licenses/LICENSE-2.0
 *
 *    Unless required by applicable law or agreed to in writing, software
 *    distributed under the License is distributed on an "AS IS" BASIS,
 *    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *    See the License for the specific language governing permissions and
 *    limitations under the License.
 */

package sqltrace

type (
	// Config represents a configuration.
	Config struct {
		Trace
	}
	// Trace represents a tracing configuration.
	Trace struct {
		Name     string  `yaml:"name"`
		Endpoint string  `yaml:"endpoint"`
		Sampler  float64 `yaml:"sampler"`
		Batcher  string  `yaml:"batcher"  validate:"eq=jaeger|eq=zipkin|eq="` // jaeger|zipkin
	}
)
