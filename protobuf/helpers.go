package rpcbench

import (
	"github.com/cgrates/cgrates/engine"
	"github.com/cgrates/cgrates/utils"
	"google/protobuf"
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
	a := &engine.AccountSummary{
		Tenant:           as.Tenant,
		ID:               as.ID,
		BalanceSummaries: make([]*engine.BalanceSummary, len(as.BalanceSummaries)), // as.BalanceSummaries.Convert(),
		AllowNegative:    as.AllowNegative,
		Disabled:         as.Disabled,
	}
	for i, v := range as.BalanceSummaries {
		a.BalanceSummaries[i] = v.Convert()
	}
	return a
}

func NewRate(r *engine.Rate) *Rate {
	return &Rate{
		GroupIntervalStart: protobuf.DurationProto(r.GroupIntervalStart),
		Value:              r.Value,
		RateIncrement:      protobuf.DurationProto(r.RateIncrement),
		RateUnit:           protobuf.DurationProto(r.RateUnit),
	}
}

func (r *Rate) Convert() *engine.Rate {
	return &engine.Rate{
		GroupIntervalStart: protobuf.Duration(r.GroupIntervalStart),
		Value:              r.Value,
		RateIncrement:      protobuf.Duration(r.RateIncrement),
		RateUnit:           protobuf.Duration(r.RateUnit),
	}
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
		Years:      make([]int32, len(rt.Years)),
		Months:     make([]int32, len(rt.Months)),
		MonthDays:  make([]int32, len(rt.MonthDays)),
		WeekDays:   make([]int32, len(rt.WeekDays)),
		StartTime:  rt.StartTime,
		EndTime:    rt.EndTime,
		cronString: rt.cronString,
		tag:        rt.tag,
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
		Years:      make(utils.Years, len(rt.Years)),
		Months:     make(utils.Months, len(rt.Months)),
		MonthDays:  make(utils.MonthDays, len(rt.MonthDays)),
		WeekDays:   make(utils.WeekDays, len(rt.WeekDays)),
		StartTime:  rt.StartTime,
		EndTime:    rt.EndTime,
		cronString: rt.cronString,
		tag:        rt.tag,
	}
	for i, v := range r.Years {
		rt.Years[i] = int(v)
	}
	for i, v := range r.Months {
		rt.Months[i] = utils.Month(v)
	}
	for i, v := range r.MonthDays {
		rt.MonthDays[i] = int(v)
	}
	for i, v := range r.WeekDays {
		rt.WeekDays[i] = utils.WeekDay(v)
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
		ActivationTime: protobuf.TimestampProto(r.ActivationTime),
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
		ActivationTime: protobuf.Timestamp(r.ActivationTime),
		RateIntervals:  make([]*engine.RateInterval, len(r.RateIntervals)),
		FallbackKeys:   make([]string, len(r.FallbackKeys)),
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
		Duration:       protobuf.TimestampProto(r.Duration),
		Cost:           r.Cost,
		BalanceInfo:    NewDebitInfo(r.BalanceInfo),
		CompressFactor: r.CompressFactor,
	}
}

func (r *Increment) Convert() *engine.Increment {
	return &engine.Increment{
		Duration:       protobuf.Timestamp(r.Duration),
		Cost:           r.Cost,
		BalanceInfo:    r.BalanceInfo.Convert(),
		CompressFactor: r.CompressFactor,
	}
}

func NewTimeSpan(r *engine.TimeSpan) *TimeSpan {
	return &TimeSpan{
		TimeStart:      protobuf.TimestampProto(r.TimeStart),
		TimeEnd:        protobuf.TimestampProto(r.TimeEnd),
		Cost:           r.Cost,
		RateInterval:   r.RateInterval.Convert(),
		DurationIndex:  protobuf.Duration(r.DurationIndex),
		Increments:     make([]*Increment, len(r.Increments)),
		RoundIncrement: NewIncrement(r.RoundIncrement),
		MatchedSubject: r.MatchedSubject,
		MatchedPrefix:  r.MatchedPrefix,
		MatchedDestId:  r.MatchedDestId,
		RatingPlanId:   r.RatingPlanId,
		CompressFactor: int32(CompressFactor),
	}
}

// func (r *TimeSpan) Convert() *engine.TimeSpan {
// 	return &engine.TimeSpan{
// 		Duration:       protobuf.Timestamp(r.Duration),
// 		Cost:           r.Cost,
// 		BalanceInfo:    r.BalanceInfo.Convert(),
// 		CompressFactor: r.CompressFactor,
// 	}
// }
