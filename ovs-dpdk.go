package main

import (
        "fmt"
        "os/exec"
        "strconv"
        "strings"
)

type Runes []rune

func IsDPDKEnabled() (bool) {
    cmd := "sudo ovs-vsctl get Open_vSwitch . dpdk_initialized | grep 'true'"
    cmdOutput, err := exec.Command("sh", "-c", cmd).Output()
    if err == nil {
        if strings.Trim(string(cmdOutput), " \n") == "true" {
            return true
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

func GetCpusFromMaskValue(maskVal string) (string) {
    intMask, err := strconv.ParseUint(maskVal, 16, 64) 
    fmt.Println(intMask)
    if err == nil {
        binMask := strings.Replace(strconv.FormatInt(intMask, 2), "0b", "", 1)
        revMask := Reverse(binMask)
        fmt.Println(binMask)
        fmt.Println(revMask)
        return revMask
    }
    return ""
}

func GetPMDCpus() (string, error) {
    cmd := "sudo ovs-vsctl --no-wait get Open_vSwitch . other_config:pmd-cpu-mask"
    cmdOutput, err := exec.Command("sh", "-c", cmd).Output()
    if err == nil {
        // fmt.Println(string(cmdOutput))
        pmdMask := GetCpusFromMaskValue(string(cmdOutput))
        return pmdMask, nil
    }
    return "", err
}

func main() {
    result := IsDPDKEnabled()
    fmt.Println(result)
    if result == true {
        fmt.Println("DPDK is intialized!")
        pmdCpus, err := GetPMDCpus()
        if err == nil {
           fmt.Println(pmdCpus)
        } 
    }
}


