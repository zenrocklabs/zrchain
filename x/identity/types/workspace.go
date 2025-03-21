package types

import (
	errorsmod "cosmossdk.io/errors"
	"github.com/Zenrock-Foundation/zrchain/v6/policy"
)

func (w *Workspace) SetAddress(addr string) { w.Address = addr }

func (w *Workspace) IsOwner(address string) bool {
	for _, owner := range w.Owners {
		if owner == address {
			return true
		}
	}
	return false
}

func (w *Workspace) IsDifferent(adminPolicyID uint64, sigPolicyID uint64) bool {
	if w.AdminPolicyId == adminPolicyID && w.SignPolicyId == sigPolicyID {
		return false
	}
	return true
}

func (w *Workspace) AddOwner(address string) error {
	if w.IsOwner(address) {
		return errorsmod.Wrapf(ErrInvalidArgument, "owner %s already is owner", address)
	}
	w.Owners = append(w.Owners, address)
	return nil
}

func (w *Workspace) RemoveOwner(address string) {
	for i, owner := range w.Owners {
		if owner == address {
			w.Owners = append(w.Owners[:i], w.Owners[i+1:]...)
			return
		}
	}
}

func (w *Workspace) AddChild(child *Workspace) {
	w.ChildWorkspaces = append(w.ChildWorkspaces, child.Address)
}

func (w *Workspace) IsChild(address string) bool {
	for _, child := range w.ChildWorkspaces {
		if child == address {
			return true
		}
	}
	return false
}

func (w *Workspace) PolicyAddOwner() policy.Policy {
	return w.AnyOwnerPolicy()
}

func (w *Workspace) PolicyRemoveOwner() policy.Policy {
	return w.AnyOwnerPolicy()
}

func (w *Workspace) PolicyAppendChild() policy.Policy {
	return w.AnyOwnerPolicy()
}

func (w *Workspace) PolicyNewKeyRequest() policy.Policy {
	return w.AnyOwnerPolicy()
}

func (w *Workspace) PolicyNewSignatureRequest() policy.Policy {
	return w.AnyOwnerPolicy()
}

func (w *Workspace) PolicyNewSignTransactionRequest() policy.Policy {
	return w.AnyOwnerPolicy()
}

func (w *Workspace) PolicyUpdateWorkspace() policy.Policy {
	return w.AnyOwnerPolicy()
}

func (w *Workspace) PolicyUpdateKeyPolicy() policy.Policy {
	return w.AnyOwnerPolicy()
}

// AnyOwnerPolicy returns a policy that is satisfied when at least one of the owners of the workspace approves.
func (w *Workspace) AnyOwnerPolicy() policy.Policy {
	return policy.NewAnyInGroupPolicy(w.Owners)
}
