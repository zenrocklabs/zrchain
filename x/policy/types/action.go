package types

import "fmt"

func (a *Action) AddApprover(approver string) error {
	if a.Status != ActionStatus_ACTION_STATUS_PENDING {
		return fmt.Errorf("action already completed")
	}

	alreadyAdded := false
	for _, a := range a.Approvers {
		if a == approver {
			alreadyAdded = true
			break
		}
	}

	if !alreadyAdded {
		a.Approvers = append(a.Approvers, approver)
	}
	return nil
}

func (a *Action) GetPolicyDataMap() map[string][]byte {
	m := make(map[string][]byte, len(a.PolicyData))
	for _, d := range a.PolicyData {
		m[d.Key] = d.Value
	}
	return m
}
