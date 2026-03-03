package utils

import "pysellers-erp-go/models"

func CalculatePayroll(s models.SalaryStructure) models.Payroll {

	gross := s.Basic + s.HRA + s.Allowance + s.Bonus

	pf := 0.0
	if s.PFEnabled {
		pf = s.Basic * 0.12 // 12% PF
	}

	esi := 0.0
	if s.ESIEnabled {
		esi = gross * 0.0075 // 0.75% ESI
	}

	tds := 0.0
	if s.TDSEnabled {
		tds = gross * 0.05 // flat 5% (can improve later)
	}

	net := gross - (pf + esi + tds)

	return models.Payroll{
		Basic:     s.Basic,
		HRA:       s.HRA,
		Allowance: s.Allowance,
		Bonus:     s.Bonus,
		Gross:     gross,
		PF:        pf,
		ESI:       esi,
		TDS:       tds,
		Net:       net,
	}
}
