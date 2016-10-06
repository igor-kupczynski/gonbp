/*
Copyright 2016 Igor Kupczynski

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
package gonbp

import (
  "io"
  "log"
  "net/http"
  "os"
)

const apiRoot = "http://api.nbp.pl/api"
const rates = apiRoot + "/exchangerates/rates/"

// Current exchange rate for a currency
func Current(table string, currency string) (string, error) {
  response, err := http.Get(rates + "/" + table + "/" + currency)
  if err != nil {
          log.Fatal(err)
          return "", err
  }
  defer response.Body.Close()
  _, err = io.Copy(os.Stdout, response.Body)
  if err != nil {
          log.Fatal(err)
          return "", err
  }
  return "", nil
}
