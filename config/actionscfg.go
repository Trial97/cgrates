/*
Real-time Online/Offline Charging System (OCS) for Telecom & ISP environments
Copyright (C) ITsysCOM GmbH

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program.  If not, see <http://www.gnu.org/licenses/>
*/

package config

import "github.com/cgrates/cgrates/utils"

// ActionSCfg is the configuration of ActionS
type ActionSCfg struct {
	Enabled             bool
	CDRsConns           []string
	EEsConns            []string
	ThresholdSConns     []string
	StatSConns          []string
	AccountSConns       []string
	Tenants             *[]string
	IndexedSelects      bool
	StringIndexedFields *[]string
	PrefixIndexedFields *[]string
	SuffixIndexedFields *[]string
	NestedFields        bool
}

func (acS *ActionSCfg) loadFromJSONCfg(jsnCfg *ActionSJsonCfg) (err error) {
	if jsnCfg == nil {
		return
	}
	if jsnCfg.Cdrs_conns != nil {
		acS.CDRsConns = make([]string, len(*jsnCfg.Cdrs_conns))
		for idx, connID := range *jsnCfg.Cdrs_conns {
			// if we have the connection internal we change the name so we can have internal rpc for each subsystem
			acS.CDRsConns[idx] = connID
			if connID == utils.MetaInternal {
				acS.CDRsConns[idx] = utils.ConcatenatedKey(utils.MetaInternal, utils.MetaCDRs)
			}
		}
	}
	if jsnCfg.Ees_conns != nil {
		acS.EEsConns = make([]string, len(*jsnCfg.Ees_conns))
		for idx, connID := range *jsnCfg.Ees_conns {
			// if we have the connection internal we change the name so we can have internal rpc for each subsystem
			acS.EEsConns[idx] = connID
			if connID == utils.MetaInternal {
				acS.EEsConns[idx] = utils.ConcatenatedKey(utils.MetaInternal, utils.MetaEEs)
			}
		}
	}
	if jsnCfg.Thresholds_conns != nil {
		acS.ThresholdSConns = make([]string, len(*jsnCfg.Thresholds_conns))
		for idx, connID := range *jsnCfg.Thresholds_conns {
			// if we have the connection internal we change the name so we can have internal rpc for each subsystem
			acS.ThresholdSConns[idx] = connID
			if connID == utils.MetaInternal {
				acS.ThresholdSConns[idx] = utils.ConcatenatedKey(utils.MetaInternal, utils.MetaThresholds)
			}
		}
	}
	if jsnCfg.Stats_conns != nil {
		acS.StatSConns = make([]string, len(*jsnCfg.Stats_conns))
		for idx, connID := range *jsnCfg.Stats_conns {
			// if we have the connection internal we change the name so we can have internal rpc for each subsystem
			acS.StatSConns[idx] = connID
			if connID == utils.MetaInternal {
				acS.StatSConns[idx] = utils.ConcatenatedKey(utils.MetaInternal, utils.MetaStats)
			}
		}
	}
	if jsnCfg.Accounts_conns != nil {
		acS.AccountSConns = make([]string, len(*jsnCfg.Accounts_conns))
		for idx, connID := range *jsnCfg.Accounts_conns {
			// if we have the connection internal we change the name so we can have internal rpc for each subsystem
			acS.AccountSConns[idx] = connID
			if connID == utils.MetaInternal {
				acS.AccountSConns[idx] = utils.ConcatenatedKey(utils.MetaInternal, utils.MetaAccounts)
			}
		}
	}
	if jsnCfg.Enabled != nil {
		acS.Enabled = *jsnCfg.Enabled
	}
	if jsnCfg.Tenants != nil {
		tnt := make([]string, len(*jsnCfg.Tenants))
		for i, fID := range *jsnCfg.Tenants {
			tnt[i] = fID
		}
		acS.Tenants = &tnt
	}
	if jsnCfg.Indexed_selects != nil {
		acS.IndexedSelects = *jsnCfg.Indexed_selects
	}
	if jsnCfg.String_indexed_fields != nil {
		sif := make([]string, len(*jsnCfg.String_indexed_fields))
		for i, fID := range *jsnCfg.String_indexed_fields {
			sif[i] = fID
		}
		acS.StringIndexedFields = &sif
	}
	if jsnCfg.Prefix_indexed_fields != nil {
		pif := make([]string, len(*jsnCfg.Prefix_indexed_fields))
		for i, fID := range *jsnCfg.Prefix_indexed_fields {
			pif[i] = fID
		}
		acS.PrefixIndexedFields = &pif
	}
	if jsnCfg.Suffix_indexed_fields != nil {
		sif := make([]string, len(*jsnCfg.Suffix_indexed_fields))
		for i, fID := range *jsnCfg.Suffix_indexed_fields {
			sif[i] = fID
		}
		acS.SuffixIndexedFields = &sif
	}
	if jsnCfg.Nested_fields != nil {
		acS.NestedFields = *jsnCfg.Nested_fields
	}
	return
}

// AsMapInterface returns the config as a map[string]interface{}
func (acS *ActionSCfg) AsMapInterface() (initialMP map[string]interface{}) {
	initialMP = map[string]interface{}{
		utils.EnabledCfg:        acS.Enabled,
		utils.IndexedSelectsCfg: acS.IndexedSelects,
		utils.NestedFieldsCfg:   acS.NestedFields,
	}
	if acS.CDRsConns != nil {
		CDRsConns := make([]string, len(acS.CDRsConns))
		for i, item := range acS.CDRsConns {
			CDRsConns[i] = item
			if item == utils.ConcatenatedKey(utils.MetaInternal, utils.MetaCDRs) {
				CDRsConns[i] = utils.MetaInternal
			}
		}
		initialMP[utils.CDRsConnsCfg] = CDRsConns
	}
	if acS.ThresholdSConns != nil {
		threshSConns := make([]string, len(acS.ThresholdSConns))
		for i, item := range acS.ThresholdSConns {
			threshSConns[i] = item
			if item == utils.ConcatenatedKey(utils.MetaInternal, utils.MetaThresholds) {
				threshSConns[i] = utils.MetaInternal
			}
		}
		initialMP[utils.ThresholdSConnsCfg] = threshSConns
	}
	if acS.StatSConns != nil {
		statSConns := make([]string, len(acS.StatSConns))
		for i, item := range acS.StatSConns {
			statSConns[i] = item
			if item == utils.ConcatenatedKey(utils.MetaInternal, utils.MetaStats) {
				statSConns[i] = utils.MetaInternal
			}
		}
		initialMP[utils.StatSConnsCfg] = statSConns
	}
	if acS.AccountSConns != nil {
		accountSConns := make([]string, len(acS.AccountSConns))
		for i, item := range acS.AccountSConns {
			accountSConns[i] = item
			if item == utils.ConcatenatedKey(utils.MetaInternal, utils.MetaAccounts) {
				accountSConns[i] = utils.MetaInternal
			}
		}
		initialMP[utils.AccountSConnsCfg] = accountSConns
	}
	if acS.EEsConns != nil {
		eesConns := make([]string, len(acS.EEsConns))
		for i, item := range acS.EEsConns {
			eesConns[i] = item
			if item == utils.ConcatenatedKey(utils.MetaInternal, utils.MetaEEs) {
				eesConns[i] = utils.MetaInternal
			}
		}
		initialMP[utils.EEsConnsCfg] = eesConns
	}
	if acS.Tenants != nil {
		Tenants := make([]string, len(*acS.Tenants))
		for i, item := range *acS.Tenants {
			Tenants[i] = item
		}
		initialMP[utils.Tenants] = Tenants
	}
	if acS.StringIndexedFields != nil {
		stringIndexedFields := make([]string, len(*acS.StringIndexedFields))
		for i, item := range *acS.StringIndexedFields {
			stringIndexedFields[i] = item
		}
		initialMP[utils.StringIndexedFieldsCfg] = stringIndexedFields
	}
	if acS.PrefixIndexedFields != nil {
		prefixIndexedFields := make([]string, len(*acS.PrefixIndexedFields))
		for i, item := range *acS.PrefixIndexedFields {
			prefixIndexedFields[i] = item
		}
		initialMP[utils.PrefixIndexedFieldsCfg] = prefixIndexedFields
	}
	if acS.SuffixIndexedFields != nil {
		suffixIndexedFields := make([]string, len(*acS.SuffixIndexedFields))
		for i, item := range *acS.SuffixIndexedFields {
			suffixIndexedFields[i] = item
		}
		initialMP[utils.SuffixIndexedFieldsCfg] = suffixIndexedFields
	}
	return
}

// Clone returns a deep copy of ActionSCfg
func (acS ActionSCfg) Clone() (cln *ActionSCfg) {
	cln = &ActionSCfg{
		Enabled:        acS.Enabled,
		IndexedSelects: acS.IndexedSelects,
		NestedFields:   acS.NestedFields,
	}
	if acS.CDRsConns != nil {
		cln.CDRsConns = make([]string, len(acS.CDRsConns))
		for i, con := range acS.CDRsConns {
			cln.CDRsConns[i] = con
		}
	}
	if acS.ThresholdSConns != nil {
		cln.ThresholdSConns = make([]string, len(acS.ThresholdSConns))
		for i, con := range acS.ThresholdSConns {
			cln.ThresholdSConns[i] = con
		}
	}
	if acS.StatSConns != nil {
		cln.StatSConns = make([]string, len(acS.StatSConns))
		for i, con := range acS.StatSConns {
			cln.StatSConns[i] = con
		}
	}
	if acS.AccountSConns != nil {
		cln.AccountSConns = make([]string, len(acS.AccountSConns))
		for i, con := range acS.AccountSConns {
			cln.AccountSConns[i] = con
		}
	}
	if acS.EEsConns != nil {
		cln.EEsConns = make([]string, len(acS.EEsConns))
		for i, k := range acS.EEsConns {
			cln.EEsConns[i] = k
		}
	}
	if acS.Tenants != nil {
		tnt := make([]string, len(*acS.Tenants))
		for i, dx := range *acS.Tenants {
			tnt[i] = dx
		}
		cln.Tenants = &tnt
	}
	if acS.StringIndexedFields != nil {
		idx := make([]string, len(*acS.StringIndexedFields))
		for i, dx := range *acS.StringIndexedFields {
			idx[i] = dx
		}
		cln.StringIndexedFields = &idx
	}
	if acS.PrefixIndexedFields != nil {
		idx := make([]string, len(*acS.PrefixIndexedFields))
		for i, dx := range *acS.PrefixIndexedFields {
			idx[i] = dx
		}
		cln.PrefixIndexedFields = &idx
	}
	if acS.SuffixIndexedFields != nil {
		idx := make([]string, len(*acS.SuffixIndexedFields))
		for i, dx := range *acS.SuffixIndexedFields {
			idx[i] = dx
		}
		cln.SuffixIndexedFields = &idx
	}
	return
}
