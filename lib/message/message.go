/*
Copyright © 2024 Kaley Main <kaleymain@google.com>

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
package message

import (
	"fmt"
	"os"

	"github.com/fatih/color"
)

// Toggles debug logging. Set to true to log debug
// level messages.
var DebugEnabled bool = false

// Prints an info level message to the output.
func Info(format string, a ...any) {
	// No color, just plain output.
	fmt.Printf(format+"\n", a...)
}

// Prints a debug level message to the output.
func Debug(format string, a ...any) {
	if DebugEnabled {
		color.Cyan(format, a...)
	}
}

// Prints a warning level message to the output.
func Warning(format string, a ...any) {
	color.Yellow(format, a...)
}

// Prints an error level message to the output,
// quitting with an exit status of 1.
func Error(format string, a ...any) {
	color.Red(format, a...)
	os.Exit(1)
}
