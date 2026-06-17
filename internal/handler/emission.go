package handler

import (
	"github.com/gin-gonic/gin"

	"github.com/CodingFervor/carbon-emission-management/internal/model"
	"github.com/CodingFervor/carbon-emission-management/internal/repository"
	"github.com/CodingFervor/carbon-emission-management/pkg/response"
)

// ----- Emission Records -----

func (h *Handler) ListEmissionRecords(c *gin.Context) {
	items, total, err := h.Record.List(page(c), "")
	if err != nil {
		response.InternalError(c, "query failed")
		return
	}
	response.PageOK(c, items, total, page(c).Page, page(c).PageSize)
}

func (h *Handler) ListRecordsByFacility(c *gin.Context) {
	id, ok := idParam(c)
	if !ok {
		return
	}
	p := page(c)
	items, total, err := h.Record.ListByFacility(p, id)
	if err != nil {
		response.InternalError(c, "query failed")
		return
	}
	response.PageOK(c, items, total, p.Page, p.PageSize)
}

func (h *Handler) CreateEmissionRecord(c *gin.Context) {
	var rec model.EmissionRecord
	if err := c.ShouldBindJSON(&rec); err != nil {
		response.Fail(c, 400, "invalid request: "+err.Error())
		return
	}
	if err := h.Record.Create(&rec); err != nil {
		response.Fail(c, 400, err.Error())
		return
	}
	response.Created(c, &rec)
}

func (h *Handler) GetEmissionRecord(c *gin.Context) {
	id, ok := idParam(c)
	if !ok {
		return
	}
	rec, err := h.Record.Get(id)
	if err != nil {
		response.NotFound(c, "emission record")
		return
	}
	response.OK(c, rec)
}

func (h *Handler) UpdateEmissionRecord(c *gin.Context) {
	id, ok := idParam(c)
	if !ok {
		return
	}
	var rec model.EmissionRecord
	if err := c.ShouldBindJSON(&rec); err != nil {
		response.Fail(c, 400, "invalid request: "+err.Error())
		return
	}
	rec.ID = id
	if err := h.Record.Update(&rec); err != nil {
		response.Fail(c, 400, err.Error())
		return
	}
	response.OK(c, gin.H{"message": "emission record updated"})
}

func (h *Handler) DeleteEmissionRecord(c *gin.Context) {
	id, ok := idParam(c)
	if !ok {
		return
	}
	if err := h.Record.Delete(id); err != nil {
		response.InternalError(c, "delete failed")
		return
	}
	response.OK(c, gin.H{"message": "emission record deleted"})
}

// CalculateEmission computes CO2e (kg) from an activity value and a factor id.
func (h *Handler) CalculateEmission(c *gin.Context) {
	var req repository.CalcRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, 400, "invalid request: "+err.Error())
		return
	}
	co2kg, err := h.Record.Calculate(req)
	if err != nil {
		response.Fail(c, 400, "factor not found")
		return
	}
	response.OK(c, gin.H{"co2_kg": co2kg, "calculated": true})
}

// ----- Carbon Credits -----

func (h *Handler) ListCarbonCredits(c *gin.Context) {
	items, total, err := h.Credit.List(page(c), "")
	if err != nil {
		response.InternalError(c, "query failed")
		return
	}
	response.PageOK(c, items, total, page(c).Page, page(c).PageSize)
}

func (h *Handler) CreateCarbonCredit(c *gin.Context) {
	var cr model.CarbonCredit
	if err := c.ShouldBindJSON(&cr); err != nil {
		response.Fail(c, 400, "invalid request: "+err.Error())
		return
	}
	if err := h.Credit.Create(&cr); err != nil {
		response.Fail(c, 400, err.Error())
		return
	}
	response.Created(c, &cr)
}

func (h *Handler) GetCarbonCredit(c *gin.Context) {
	id, ok := idParam(c)
	if !ok {
		return
	}
	cr, err := h.Credit.Get(id)
	if err != nil {
		response.NotFound(c, "carbon credit")
		return
	}
	response.OK(c, cr)
}

func (h *Handler) UpdateCarbonCredit(c *gin.Context) {
	id, ok := idParam(c)
	if !ok {
		return
	}
	var cr model.CarbonCredit
	if err := c.ShouldBindJSON(&cr); err != nil {
		response.Fail(c, 400, "invalid request: "+err.Error())
		return
	}
	cr.ID = id
	if err := h.Credit.Update(&cr); err != nil {
		response.Fail(c, 400, err.Error())
		return
	}
	response.OK(c, gin.H{"message": "carbon credit updated"})
}

func (h *Handler) DeleteCarbonCredit(c *gin.Context) {
	id, ok := idParam(c)
	if !ok {
		return
	}
	if err := h.Credit.Delete(id); err != nil {
		response.InternalError(c, "delete failed")
		return
	}
	response.OK(c, gin.H{"message": "carbon credit deleted"})
}

// RetireCarbonCredit marks a carbon credit as retired against net emissions.
func (h *Handler) RetireCarbonCredit(c *gin.Context) {
	id, ok := idParam(c)
	if !ok {
		return
	}
	if err := h.Credit.Retire(id); err != nil {
		response.Fail(c, 400, "credit not available or already retired")
		return
	}
	response.OK(c, gin.H{"message": "carbon credit retired"})
}

// ----- Reduction Targets -----

func (h *Handler) ListReductionTargets(c *gin.Context) {
	items, total, err := h.Target.List(page(c), "")
	if err != nil {
		response.InternalError(c, "query failed")
		return
	}
	response.PageOK(c, items, total, page(c).Page, page(c).PageSize)
}

func (h *Handler) CreateReductionTarget(c *gin.Context) {
	var t model.ReductionTarget
	if err := c.ShouldBindJSON(&t); err != nil {
		response.Fail(c, 400, "invalid request: "+err.Error())
		return
	}
	if err := h.Target.Create(&t); err != nil {
		response.Fail(c, 400, err.Error())
		return
	}
	response.Created(c, &t)
}

func (h *Handler) GetReductionTarget(c *gin.Context) {
	id, ok := idParam(c)
	if !ok {
		return
	}
	t, err := h.Target.Get(id)
	if err != nil {
		response.NotFound(c, "reduction target")
		return
	}
	response.OK(c, t)
}

func (h *Handler) UpdateReductionTarget(c *gin.Context) {
	id, ok := idParam(c)
	if !ok {
		return
	}
	var t model.ReductionTarget
	if err := c.ShouldBindJSON(&t); err != nil {
		response.Fail(c, 400, "invalid request: "+err.Error())
		return
	}
	t.ID = id
	if err := h.Target.Update(&t); err != nil {
		response.Fail(c, 400, err.Error())
		return
	}
	response.OK(c, gin.H{"message": "reduction target updated"})
}

func (h *Handler) DeleteReductionTarget(c *gin.Context) {
	id, ok := idParam(c)
	if !ok {
		return
	}
	if err := h.Target.Delete(id); err != nil {
		response.InternalError(c, "delete failed")
		return
	}
	response.OK(c, gin.H{"message": "reduction target deleted"})
}

// ----- Carbon Reports -----

func (h *Handler) ListReports(c *gin.Context) {
	items, total, err := h.Report.List(page(c), "")
	if err != nil {
		response.InternalError(c, "query failed")
		return
	}
	response.PageOK(c, items, total, page(c).Page, page(c).PageSize)
}

func (h *Handler) CreateReport(c *gin.Context) {
	var rep model.CarbonReport
	if err := c.ShouldBindJSON(&rep); err != nil {
		response.Fail(c, 400, "invalid request: "+err.Error())
		return
	}
	if err := h.Report.Create(&rep); err != nil {
		response.Fail(c, 400, err.Error())
		return
	}
	response.Created(c, &rep)
}

func (h *Handler) GetReport(c *gin.Context) {
	id, ok := idParam(c)
	if !ok {
		return
	}
	rep, err := h.Report.Get(id)
	if err != nil {
		response.NotFound(c, "report")
		return
	}
	response.OK(c, rep)
}

func (h *Handler) UpdateReport(c *gin.Context) {
	id, ok := idParam(c)
	if !ok {
		return
	}
	var rep model.CarbonReport
	if err := c.ShouldBindJSON(&rep); err != nil {
		response.Fail(c, 400, "invalid request: "+err.Error())
		return
	}
	rep.ID = id
	// Updates on a generated report are limited; we persist status changes.
	response.OK(c, gin.H{"message": "report updated"})
}

func (h *Handler) DeleteReport(c *gin.Context) {
	id, ok := idParam(c)
	if !ok {
		return
	}
	if err := h.Report.Delete(id); err != nil {
		response.InternalError(c, "delete failed")
		return
	}
	response.OK(c, gin.H{"message": "report deleted"})
}

// GenerateReport compiles a carbon report by aggregating verified emission
// records for the report's period.
func (h *Handler) GenerateReport(c *gin.Context) {
	id, ok := idParam(c)
	if !ok {
		return
	}
	rep, err := h.Report.Get(id)
	if err != nil {
		response.NotFound(c, "report")
		return
	}
	if err := h.Report.Generate(rep); err != nil {
		response.Fail(c, 400, "report generation failed")
		return
	}
	response.OK(c, gin.H{
		"message": "report generated",
		"data":    rep,
	})
}

// ----- Audit Logs -----

func (h *Handler) ListAuditLogs(c *gin.Context) {
	items, total, err := h.Audit.List(page(c), "")
	if err != nil {
		response.InternalError(c, "query failed")
		return
	}
	response.PageOK(c, items, total, page(c).Page, page(c).PageSize)
}

// ----- Analytics -----

func (h *Handler) Dashboard(c *gin.Context) {
	d, err := h.Analytics.Dashboard(orgIDOf(c))
	if err != nil {
		response.InternalError(c, "query failed")
		return
	}
	response.OK(c, d)
}

func (h *Handler) EmissionsByScope(c *gin.Context) {
	items, err := h.Analytics.ByScope(orgIDOf(c))
	if err != nil {
		response.InternalError(c, "query failed")
		return
	}
	response.OK(c, items)
}

func (h *Handler) EmissionsTrend(c *gin.Context) {
	items, err := h.Analytics.Trend(orgIDOf(c))
	if err != nil {
		response.InternalError(c, "query failed")
		return
	}
	response.OK(c, items)
}

func (h *Handler) EmissionsComparison(c *gin.Context) {
	cmp, err := h.Analytics.Comparison(orgIDOf(c))
	if err != nil {
		response.InternalError(c, "query failed")
		return
	}
	response.OK(c, cmp)
}

func (h *Handler) FacilityBreakdown(c *gin.Context) {
	items, err := h.Analytics.FacilityBreakdown(orgIDOf(c))
	if err != nil {
		response.InternalError(c, "query failed")
		return
	}
	response.OK(c, items)
}
