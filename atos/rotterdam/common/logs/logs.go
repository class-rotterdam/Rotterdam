//
// Copyright 2018 Atos
//
// ROTTERDAM application
// CLASS Project: https://class-project.eu/
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//     https://www.apache.org/licenses/LICENSE-2.0
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//
// @author: ATOS
//

package logs

import (
	log "github.com/sirupsen/logrus"
)

func init() {
	log.SetFormatter(&log.TextFormatter{ForceColors: true, FullTimestamp: true})
	//log.SetFormatter(&log.JSONFormatter{})
	log.SetLevel(log.DebugLevel)
}

///////////////////////////////////////////////////////////////////////////////
// Println, Printf functions ==> INFO

/*

 */
func Println(m string) {
	log.Println(m)
}

/*

 */
func Printlne(m string, e error) {
	log.Error(m, e)
}

/*


func Printf(m string) {
	log.Printf(m)
} */

/*

 */
func Printf(format string, args ...interface{}) {
	log.Printf(format, args...)
}

/*

 */
func Printfi(m string, s int) {
	log.Printf(m, s)
}

/*

 */
func Fatal(e error) {
	log.Fatal(e)
}

///////////////////////////////////////////////////////////////////////////////
// Logger functions

/*

 */
func Trace(m string) {
	log.Trace(m)
}

/*

 */
func Debug(m string) {
	log.Debug(m)
}

/*

 */
func Info(m string) {
	log.Info(m)
}

/*

 */
func Warn(m string) {
	log.Warn(m)
}

/*

 */
func Error(args ...interface{}) {
	log.Error(args...)
}

/*

 */
func Panic(m string) {
	log.Panic(m)
}
