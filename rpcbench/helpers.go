package rpcbench

import (
	"time"

	"github.com/cgrates/cgrates/engine"
	"github.com/cgrates/cgrates/utils"
)

func NewBalanceSummary(bs *engine.BalanceSummary) *BalanceSummary {
	return &BalanceSummary{
		UUID:     bs.UUID,
		ID:       bs.ID,
		Type:     bs.Type,
		Value:    bs.Value,
		Disabled: bs.Disabled,
	}
}
func (bs *BalanceSummary) Convert() *engine.BalanceSummary {
	return &engine.BalanceSummary{
		UUID:     bs.UUID,
		ID:       bs.ID,
		Type:     bs.Type,
		Value:    bs.Value,
		Disabled: bs.Disabled,
	}
}

func NewAccountSummary(as *engine.AccountSummary) *AccountSummary {
	a := &AccountSummary{
		Tenant:           as.Tenant,
		ID:               as.ID,
		BalanceSummaries: make([]*BalanceSummary, len(as.BalanceSummaries)), //NewBalanceSummary(as.BalanceSummaries),
		AllowNegative:    as.AllowNegative,
		Disabled:         as.Disabled,
	}
	for i, v := range as.BalanceSummaries {
		a.BalanceSummaries[i] = NewBalanceSummary(v)
	}
	return a
}

func (as *AccountSummary) Convert() *engine.AccountSummary {
	if as == nil {
		return nil
	}
	a := &engine.AccountSummary{
		Tenant:        as.Tenant,
		ID:            as.ID,
		AllowNegative: as.AllowNegative,
		Disabled:      as.Disabled,
	}
	if as.BalanceSummaries != nil {
		a.BalanceSummaries = make([]*engine.BalanceSummary, len(as.BalanceSummaries))
		for i, v := range as.BalanceSummaries {
			a.BalanceSummaries[i] = v.Convert()
		}
	}
	return a
}

func NewRate(r *engine.Rate) *Rate {
	return &Rate{
		GroupIntervalStart: &r.GroupIntervalStart,
		Value:              r.Value,
		RateIncrement:      &r.RateIncrement,
		RateUnit:           &r.RateUnit,
	}
}

func (r *Rate) Convert() *engine.Rate {
	rt := &engine.Rate{
		// GroupIntervalStart: r.GroupIntervalStart,
		Value: r.Value,
		// RateIncrement:      r.RateIncrement,
		// RateUnit:           r.RateUnit,
	}
	if r.GroupIntervalStart != nil {
		rt.GroupIntervalStart = *r.GroupIntervalStart
	}
	if r.RateIncrement != nil {
		rt.RateIncrement = *r.RateIncrement
	}
	if r.RateUnit != nil {
		rt.RateUnit = *r.RateUnit
	}
	return rt
}

func NewRIRate(r *engine.RIRate) *RIRate {
	rt := &RIRate{
		ConnectFee:       r.ConnectFee,
		RoundingMethod:   r.RoundingMethod,
		RoundingDecimals: int32(r.RoundingDecimals),
		MaxCost:          r.MaxCost,
		MaxCostStrategy:  r.MaxCostStrategy,
		Rates:            make([]*Rate, len(r.Rates)),
	}
	for i, v := range r.Rates {
		rt.Rates[i] = NewRate(v)
	}
	return rt
}

func (r *RIRate) Convert() *engine.RIRate {
	rt := &engine.RIRate{
		ConnectFee:       r.ConnectFee,
		RoundingMethod:   r.RoundingMethod,
		RoundingDecimals: int(r.RoundingDecimals),
		MaxCost:          r.MaxCost,
		MaxCostStrategy:  r.MaxCostStrategy,
		Rates:            make([]*engine.Rate, len(r.Rates)),
	}
	for i, v := range r.Rates {
		rt.Rates[i] = v.Convert()
	}
	return rt
}

func NewRITiming(r *engine.RITiming) *RITiming {
	rt := &RITiming{
		Years:     make([]int32, len(r.Years)),
		Months:    make([]int32, len(r.Months)),
		MonthDays: make([]int32, len(r.MonthDays)),
		WeekDays:  make([]int32, len(r.WeekDays)),
		StartTime: r.StartTime,
		EndTime:   r.EndTime,
	}
	for i, v := range r.Years {
		rt.Years[i] = int32(v)
	}
	for i, v := range r.Months {
		rt.Months[i] = int32(v)
	}
	for i, v := range r.MonthDays {
		rt.MonthDays[i] = int32(v)
	}
	for i, v := range r.WeekDays {
		rt.WeekDays[i] = int32(v)
	}
	return rt
}

func (r *RITiming) Convert() *engine.RITiming {
	rt := &engine.RITiming{
		Years:     make(utils.Years, len(r.Years)),
		Months:    make(utils.Months, len(r.Months)),
		MonthDays: make(utils.MonthDays, len(r.MonthDays)),
		WeekDays:  make(utils.WeekDays, len(r.WeekDays)),
		StartTime: r.StartTime,
		EndTime:   r.EndTime,
	}
	for i, v := range r.Years {
		rt.Years[i] = int(v)
	}
	for i, v := range r.Months {
		rt.Months[i] = time.Month(v)
	}
	for i, v := range r.MonthDays {
		rt.MonthDays[i] = int(v)
	}
	for i, v := range r.WeekDays {
		rt.WeekDays[i] = time.Weekday(v)
	}
	return rt
}

func NewRateInterval(r *engine.RateInterval) *RateInterval {
	return &RateInterval{
		Timing: NewRITiming(r.Timing),
		Rating: NewRIRate(r.Rating),
		Weight: r.Weight,
	}
}

func (r *RateInterval) Convert() *engine.RateInterval {
	return &engine.RateInterval{
		Timing: r.Timing.Convert(),
		Rating: r.Rating.Convert(),
		Weight: r.Weight,
	}
}

func NewRatingInfo(r *engine.RatingInfo) *RatingInfo {
	rt := &RatingInfo{
		MatchedSubject: r.MatchedSubject,
		RatingPlanId:   r.RatingPlanId,
		MatchedPrefix:  r.MatchedPrefix,
		MatchedDestId:  r.MatchedDestId,
		ActivationTime: &r.ActivationTime,
		RateIntervals:  make([]*RateInterval, len(r.RateIntervals)),
		FallbackKeys:   make([]string, len(r.FallbackKeys)),
	}
	for i, v := range r.RateIntervals {
		rt.RateIntervals[i] = NewRateInterval(v)
	}
	for i, v := range r.FallbackKeys {
		rt.FallbackKeys[i] = v
	}
	return rt
}

func (r *RatingInfo) Convert() *engine.RatingInfo {
	rt := &engine.RatingInfo{
		MatchedSubject: r.MatchedSubject,
		RatingPlanId:   r.RatingPlanId,
		MatchedPrefix:  r.MatchedPrefix,
		MatchedDestId:  r.MatchedDestId,
		// ActivationTime: r.ActivationTime,
		RateIntervals: make([]*engine.RateInterval, len(r.RateIntervals)),
		FallbackKeys:  make([]string, len(r.FallbackKeys)),
	}
	if r.ActivationTime != nil {
		r.ActivationTime = r.ActivationTime
	}
	for i, v := range r.RateIntervals {
		rt.RateIntervals[i] = v.Convert()
	}
	for i, v := range r.FallbackKeys {
		rt.FallbackKeys[i] = v
	}
	return rt
}

func NewUnitInfo(r *engine.UnitInfo) *UnitInfo {
	return &UnitInfo{
		UUID:          r.UUID,
		ID:            r.ID,
		Value:         r.Value,
		DestinationID: r.DestinationID,
		Consumed:      r.Consumed,
		TOR:           r.TOR,
		RateInterval:  NewRateInterval(r.RateInterval),
	}
}

func (r *UnitInfo) Convert() *engine.UnitInfo {
	return &engine.UnitInfo{
		UUID:          r.UUID,
		ID:            r.ID,
		Value:         r.Value,
		DestinationID: r.DestinationID,
		Consumed:      r.Consumed,
		TOR:           r.TOR,
		RateInterval:  r.RateInterval.Convert(),
	}
}

func NewMonetaryInfo(r *engine.MonetaryInfo) *MonetaryInfo {
	return &MonetaryInfo{
		UUID:         r.UUID,
		ID:           r.ID,
		Value:        r.Value,
		RateInterval: NewRateInterval(r.RateInterval),
	}
}

func (r *MonetaryInfo) Convert() *engine.MonetaryInfo {
	return &engine.MonetaryInfo{
		UUID:         r.UUID,
		ID:           r.ID,
		Value:        r.Value,
		RateInterval: r.RateInterval.Convert(),
	}
}

func NewDebitInfo(r *engine.DebitInfo) *DebitInfo {
	return &DebitInfo{
		Unit:      NewUnitInfo(r.Unit),
		Monetary:  NewMonetaryInfo(r.Monetary),
		AccountID: r.AccountID,
	}
}

func (r *DebitInfo) Convert() *engine.DebitInfo {
	return &engine.DebitInfo{
		Unit:      r.Unit.Convert(),
		Monetary:  r.Monetary.Convert(),
		AccountID: r.AccountID,
	}
}

func NewIncrement(r *engine.Increment) *Increment {
	return &Increment{
		Duration:       &r.Duration,
		Cost:           r.Cost,
		BalanceInfo:    NewDebitInfo(r.BalanceInfo),
		CompressFactor: int32(r.CompressFactor),
	}
}

func (r *Increment) Convert() *engine.Increment {
	rt := &engine.Increment{
		Cost:           r.Cost,
		BalanceInfo:    r.BalanceInfo.Convert(),
		CompressFactor: int(r.CompressFactor),
	}
	if r.Duration != nil {
		rt.Duration = *r.Duration
	}
	return rt
}

func NewTimeSpan(r *engine.TimeSpan) *TimeSpan {
	rt := &TimeSpan{
		TimeStart:      &r.TimeStart,
		TimeEnd:        &r.TimeEnd,
		Cost:           r.Cost,
		RateInterval:   NewRateInterval(r.RateInterval),
		DurationIndex:  &r.DurationIndex,
		Increments:     make([]*Increment, len(r.Increments)),
		RoundIncrement: NewIncrement(r.RoundIncrement),
		MatchedSubject: r.MatchedSubject,
		MatchedPrefix:  r.MatchedPrefix,
		MatchedDestId:  r.MatchedDestId,
		RatingPlanId:   r.RatingPlanId,
		CompressFactor: int32(r.CompressFactor),
	}
	for i, v := range r.Increments {
		rt.Increments[i] = NewIncrement(v)
	}
	return rt
}

func (r *TimeSpan) Convert() *engine.TimeSpan {
	rt := &engine.TimeSpan{
		// TimeStart:      r.TimeStart,
		// TimeEnd:        r.TimeEnd,
		Cost:         r.Cost,
		RateInterval: r.RateInterval.Convert(),
		// DurationIndex:  r.DurationIndex,
		Increments:     make([]*engine.Increment, len(r.Increments)),
		RoundIncrement: r.RoundIncrement.Convert(),
		MatchedSubject: r.MatchedSubject,
		MatchedPrefix:  r.MatchedPrefix,
		MatchedDestId:  r.MatchedDestId,
		RatingPlanId:   r.RatingPlanId,
		CompressFactor: int(r.CompressFactor),
	}
	if r.TimeStart != nil {
		rt.TimeStart = *r.TimeStart
	}
	if r.TimeEnd != nil {
		rt.TimeEnd = *r.TimeEnd
	}
	if r.DurationIndex != nil {
		rt.DurationIndex = *r.DurationIndex
	}
	for i, v := range r.Increments {
		rt.Increments[i] = v.Convert()
	}
	return rt
}

func NewCallCost(r *engine.CallCost) *CallCost {
	rt := &CallCost{
		Category:       r.Category,
		Tenant:         r.Tenant,
		Subject:        r.Subject,
		Account:        r.Account,
		Destination:    r.Destination,
		TOR:            r.TOR,
		Cost:           r.Cost,
		RatedUsage:     r.RatedUsage,
		AccountSummary: NewAccountSummary(r.AccountSummary),
		Timespans:      make([]*TimeSpan, len(r.Timespans)),
	}
	for i, v := range r.Timespans {
		rt.Timespans[i] = NewTimeSpan(v)
	}
	return rt
}

func (r *CallCost) Convert() *engine.CallCost {
	rt := &engine.CallCost{
		Category:       r.Category,
		Tenant:         r.Tenant,
		Subject:        r.Subject,
		Account:        r.Account,
		Destination:    r.Destination,
		TOR:            r.TOR,
		Cost:           r.Cost,
		RatedUsage:     r.RatedUsage,
		AccountSummary: r.AccountSummary.Convert(),
		Timespans:      make([]*engine.TimeSpan, len(r.Timespans)),
	}
	for i, v := range r.Timespans {
		rt.Timespans[i] = v.Convert()
	}
	return rt
}

func NewCallDescriptor(r *engine.CallDescriptor) *CallDescriptor {
	rt := &CallDescriptor{
		Category:            r.Category,
		Tenant:              r.Tenant,
		Subject:             r.Subject,
		Account:             r.Account,
		Destination:         r.Destination,
		LoopIndex:           r.LoopIndex,
		FallbackSubject:     r.FallbackSubject,
		TOR:                 r.TOR,
		MaxRate:             r.MaxRate,
		MaxCostSoFar:        r.MaxCostSoFar,
		CgrID:               r.CgrID,
		RunID:               r.RunID,
		ForceDuration:       r.ForceDuration,
		PerformRounding:     r.PerformRounding,
		DryRun:              r.DryRun,
		DenyNegativeAccount: r.DenyNegativeAccount,
		ExtraFields:         r.ExtraFields,
		TimeStart:           &r.TimeStart,
		TimeEnd:             &r.TimeEnd,
		DurationIndex:       &r.DurationIndex,
		MaxRateUnit:         &r.MaxRateUnit,
		RatingInfos:         make([]*RatingInfo, len(r.RatingInfos)),
		Increments:          make([]*Increment, len(r.Increments)),
	}
	for i, v := range r.RatingInfos {
		rt.RatingInfos[i] = NewRatingInfo(v)
	}
	for i, v := range r.Increments {
		rt.Increments[i] = NewIncrement(v)
	}
	return rt
}

func (r *CallDescriptor) Convert() *engine.CallDescriptor {
	rt := &engine.CallDescriptor{
		Category:            r.Category,
		Tenant:              r.Tenant,
		Subject:             r.Subject,
		Account:             r.Account,
		Destination:         r.Destination,
		LoopIndex:           r.LoopIndex,
		FallbackSubject:     r.FallbackSubject,
		TOR:                 r.TOR,
		MaxRate:             r.MaxRate,
		MaxCostSoFar:        r.MaxCostSoFar,
		CgrID:               r.CgrID,
		RunID:               r.RunID,
		ForceDuration:       r.ForceDuration,
		PerformRounding:     r.PerformRounding,
		DryRun:              r.DryRun,
		DenyNegativeAccount: r.DenyNegativeAccount,
		ExtraFields:         r.ExtraFields,
		// TimeStart:           r.TimeStart,
		// TimeEnd:       r.TimeEnd,
		// DurationIndex: r.DurationIndex,
		// MaxRateUnit:   r.MaxRateUnit,
		RatingInfos: make([]*engine.RatingInfo, len(r.RatingInfos)),
		Increments:  make([]*engine.Increment, len(r.Increments)),
	}
	if r.TimeStart != nil {
		rt.TimeStart = *r.TimeStart
	}
	if r.TimeEnd != nil {
		rt.TimeEnd = *r.TimeEnd
	}
	if r.DurationIndex != nil {
		rt.DurationIndex = *r.DurationIndex
	}
	if r.MaxRateUnit != nil {
		rt.MaxRateUnit = *r.MaxRateUnit
	}
	for i, v := range r.RatingInfos {
		rt.RatingInfos[i] = v.Convert()
	}
	for i, v := range r.Increments {
		rt.Increments[i] = v.Convert()
	}
	return rt
}
