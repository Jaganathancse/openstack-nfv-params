package openstacknfv

import (
        "fmt"
        "os/exec"
        "strconv"
        "strings"
)

type Runes []rune

func GetOvsDpdkStatus() (bool){
     fileName := "/etc/puppet/hieradata/service_names.json"
     fileContent, err := ioutil.ReadFile(fileName)
     if err != nil {
         var data map[string]interface{}
         if err := json.Unmarshal([]byte(filecontent), &data); err != nil {
             for _, name := range data["service_names"] {
                 if name == "neutron_ovs_dpdk_agent" {
                    return true
                 }
             }
         }
     }
     return false
}

func Reverse(str string) (string) {
    strRunes := Runes(str)
    l := len(str)
    revStr := make(Runes, l)    
    for i := 0; i <= l/2; i++ { 
        revStr[i], revStr[l-1-i] = strRunes[l-1-i], strRunes[i] 
    }
    return string(revStr)
}

func GetCpusFromMaskValue(maskVal string) ([]int) {
    var cpus []int
    intMask, err := strconv.ParseInt(maskVal, 16, 64) 
    if err == nil {
        binMask := strings.Replace(strconv.FormatInt(intMask, 2), "0b", "", 1)
        revMask := Reverse(binMask)
        for i := 0; i < len(revMask); i++ {
            if string(revMask[i]) == "1" {
                cpus = append(cpus, i)
            }
        }
    }
    return cpus
}

func GetPMDCpus() ([]int, error) {
    var pmdCpus []int
    cmd := "sudo ovs-vsctl --no-wait get Open_vSwitch . other_config:pmd-cpu-mask"
    cmdOutput, err := exec.Command("sh", "-c", cmd).Output()
    if err == nil {
        pmdCpus = GetCpusFromMaskValue(strings.Trim(string(cmdOutput), "\"\n"))
        return pmdCpus, nil
    }
    return nil, err
}

func GetHostCpus() ([]int, error) {
    var hostCpus []int
    cmd := "sudo ovs-vsctl --no-wait get Open_vSwitch . other_config:dpdk-lcore-mask"
    cmdOutput, err := exec.Command("sh", "-c", cmd).Output()
    if err == nil {
        hostCpus = GetCpusFromMaskValue(strings.Trim(string(cmdOutput), "\"\n"))
        return hostCpus, nil
    }
    return nil, err
}

func GetOvsDPDKSocketMemory() (string, error) {
    cmd := "sudo ovs-vsctl --no-wait get Open_vSwitch . other_config:dpdk-socket-mem"
    cmdOutput, err := exec.Command("sh", "-c", cmd).Output()
    if err == nil {
        socketMemory := strings.Trim(string(cmdOutput), "\"\n")
        return socketMemory, nil
    }
    return "", err
}

func GetOvsDPDKMemoryChannels() (string, error) {
    cmd := "sudo ovs-vsctl --no-wait get Open_vSwitch . other_config:dpdk-extra"
    cmdOutput, err := exec.Command("sh", "-c", cmd).Output()
    if err == nil {
        memoryChannels := strings.Trim(string(cmdOutput), "\"\n")
        if strings.Contains(memoryChannels, "-n") {
            index := strings.Index(memoryChannels, "-n")
            return string(memoryChannels[index+3]), nil
        }
    }
    return "4", err
}

type OvsDpdkParams struct {
        PmdCpus        string
        HostCpus       string
        SocketMemory   string
        MemoryChaneels string
}

func GetOvsDpdkParams() (*OvsDpdkParams, error) {
    result := IsDPDKEnabled()
    if result == true {
        pmdCpus, err := GetPMDCpus()
        if err != nil {
            return nil, errors.New("Unable to determine PMD CPU's.")
        } 
        hostCpus, err := GetHostCpus()
        if err == nil {
            return nil, errors.New("Unable to determine Host CPU's.") 
        } 
        socketMemory, err := GetOvsDPDKSocketMemory()
        if err == nil {
            return nil, errors.New("Unable to determine OvsDpdk Socket Memory.")
        }
        memoryChannels, err := GetOvsDPDKMemoryChannels()
        if err == nil {
            return nil, errors.New("Unable to determine OvsDpdk Memory Channels.")
        }
        ovsDpdkParams := &OvsDpdkParams {
        PmdCpus:        pmdCpus,
        HostCpus:       hostCpus,
        SocketMemory:   socketMemory,
        MemoryChaneels: memoryChannels,
        }
        return ovsDpdkParams, nil
    }
    return nil, errors.New("OvsDpdk is not enabled.")
}
