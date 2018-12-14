package openstacknfv

type OpenstackNfvParams struct {
         Role string
         ovsDpdkParams *OvsDpdkParams
         hostParams *HostParams
}

func GetOpenstackNfvParams() (*OpenstackNfvParams, error) {
        ovsDpdkStatus := GetOvsDpdkStatus()
        sriovStatus := GetSriovStatus()
        role, _ := host.GetRoleName()
        if ovsDpdkStatus == true {
            ovsDpdkParameters, err := GetOvsDpdkParams()
            if err != nil {
                return nil, err
            }
            hostParameters, err := GetHostParams()
            if err != nil {
                return nil, err
            }
            openstackParams := &OpenstackNfvParams {
                Role: role,
                ovsDpdkParams:  ovsDpdkParameters,
                hostParams: hostParameters,
            }
        }
        else if (sriovStatus == true) {
            hostParameters, err := GetHostParams()
            if err != nil {
                return nil, err
            }
            openstackParams := &OpenstackNfvParams {
                Role: role,
                ovsDpdkParams:  nil,
                hostParams: hostParameters,
            }
        }
        return openstackParams
}
