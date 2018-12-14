package openstacknfv

import (
        "os/exec"
        "strings"
)

func GetRoleName() (string, error) {
    cmd := "sudo cat /etc/role.conf | grep \"Role\" | grep -v ^#"
    cmdOutput, err := exec.Command("sh", "-c", cmd).Output()
    if err == nil {
        role := strings.Trim(string(cmdOutput), "\"\n")
        role = strings.Trim(role, "Role=")
        return role, nil
    }
    return "", err
}
func GetNovaReservedHostMemory() (string, error) {
    cmd := "sudo cat /var/lib/config-data/nova_libvirt/etc/nova/nova.conf | grep \"reserved_host_memory_mb\" | grep -v ^#"
    cmdOutput, err := exec.Command("sh", "-c", cmd).Output()
    if err == nil {
        novaReservedHostMemory := strings.Trim(string(cmdOutput), "\"\n")
        novaReservedHostMemory = strings.Trim(novaReservedHostMemory, "reserved_host_memory_mb=")
        return novaReservedHostMemory, nil
    }
    return "", err
}

func GetNovaCpus() (string, error) {
    cmd := "sudo cat /var/lib/config-data/nova_libvirt/etc/nova/nova.conf | grep \"vcpu_pin_set\" | grep -v ^#"
    cmdOutput, err := exec.Command("sh", "-c", cmd).Output()
    if err == nil {
        novaCpus := strings.Trim(string(cmdOutput), "\"\n")
        novaCpus = strings.Trim(novaCpus, "vcpu_pin_set=")
        return novaCpus, nil
    }
    return "", err
}

func GetHostIsolatedCpus() (string, error) {
    cmd := "sudo cat /etc/tuned/cpu-partitioning-variables.conf | grep \"isolated_cores\" | grep -v ^#"
    cmdOutput, err := exec.Command("sh", "-c", cmd).Output()
    if err == nil {
        hostIsolCpus := strings.Trim(string(cmdOutput), "\"\n")
        hostIsolCpus = strings.Trim(hostIsolCpus, "isolated_cores=")
        return hostIsolCpus, nil
    }
    return "", err
}

func GetKernelArgs() (string, error) {
    var kernelArgs []string
    cmd := "sudo cat /etc/default/grub | grep "TRIPLEO_HEAT_TEMPLATE_KERNEL_ARGS" | grep -v ^#'"
    cmdOutput, err := exec.Command("sh", "-c", cmd).Output()
    if err == nil {
        output := strings.Trim(string(cmdOutput), "\"\n")
        params = strings.Split(output, " ")
        if len(params) > 0 {
           for _, param := range params {
               param = strings.TrimSpace(param)
               if (strings.Index(param, "hugepages") == 0 ||
                  strings.Index(param, "intel_iommu") == 0 ||
                  strings.Index(param, "iommu") == 0) {
                  kernelArgs = append(kernelArgs, param)
               }
           }
        }
        return strings.join(kernelArgs, ","), nil
    }
    return "", err
}

type HostParams struct {
        NovaReservedMemory int64
        NovaCpus           string
        IsolCpus           string
        KernelArgs         string
}

func GetHostParams() (*HostParams, error){
     novaReservedMemory, err := GetNovaReservedHostMemory()
     if err != nil {
          return nil, errors.New("Unable to determine Nova Reserved Host Memory.")
     }
     novaCpus, err := GetNovaCpus()
     if err != nil {
          return nil, errors.New("Unable to determine Nova CPUs.")
     }
     isolCpus, err := GetHostIsolatedCpus()
     if err != nil {
          return nil, errors.New("Unable to determine Isol CPUs.")
     }
     kernelArgs, err := GetKernelArgs()
     if err != nil {
          return nil, errors.New("Unable to determine Kernel Args.")
     }
     hostParams := &HostParams {
                   NovaReservedMemory: int64(novaReservedMemory),
                   NovaCpus:           novaCpus,
                   IsolCpus:           isolCpus,
                   KernelArgs:         kernelArgs,
     }
     return hostParams, nil
}
