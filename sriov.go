package openstacknfv

import (
       "io/ioutil"
)

func GetSriovStatus() (bool){
     fileName := "/etc/puppet/hieradata/service_names.json"
     fileContent, err := ioutil.ReadFile(fileName)
     if err != nil {
         var data map[string]interface{}
         if err := json.Unmarshal([]byte(filecontent), &data); err != nil {
             for _, name := range data["service_names"] {
                 if name == "neutron_sriov_agent" {
                    return true
                 }
             }
         }
     }
     return false
}
