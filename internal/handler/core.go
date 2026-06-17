package handler

import (
	"github.com/gin-gonic/gin"

	"github.com/CodingFervor/carbon-emission-management/internal/model"
	"github.com/CodingFervor/carbon-emission-management/pkg/response"
)

// ----- Organizations -----

func (h *Handler) ListOrganizations(c *gin.Context) {
	items, total, err := h.Org.List(page(c), "status = 'active'")
	if err != nil {
		response.InternalError(c, "query failed")
		return
	}
	response.PageOK(c, items, total, page(c).Page, page(c).PageSize)
}

func (h *Handler) CreateOrganization(c *gin.Context) {
	var o model.Organization
	if err := c.ShouldBindJSON(&o); err != nil {
		response.Fail(c, 400, "invalid request: "+err.Error())
		return
	}
	if err := h.Org.Create(&o); err != nil {
		response.Fail(c, 400, err.Error())
		return
	}
	response.Created(c, &o)
}

func (h *Handler) GetOrganization(c *gin.Context) {
	id, ok := idParam(c)
	if !ok {
		return
	}
	o, err := h.Org.Get(id)
	if err != nil {
		response.NotFound(c, "organization")
		return
	}
	response.OK(c, o)
}

func (h *Handler) UpdateOrganization(c *gin.Context) {
	id, ok := idParam(c)
	if !ok {
		return
	}
	var o model.Organization
	if err := c.ShouldBindJSON(&o); err != nil {
		response.Fail(c, 400, "invalid request: "+err.Error())
		return
	}
	o.ID = id
	if err := h.Org.Update(&o); err != nil {
		response.Fail(c, 400, err.Error())
		return
	}
	response.OK(c, gin.H{"message": "organization updated"})
}

func (h *Handler) DeleteOrganization(c *gin.Context) {
	id, ok := idParam(c)
	if !ok {
		return
	}
	if err := h.Org.Delete(id); err != nil {
		response.InternalError(c, "delete failed")
		return
	}
	response.OK(c, gin.H{"message": "organization deleted"})
}

// ----- Facilities -----

func (h *Handler) ListFacilities(c *gin.Context) {
	items, total, err := h.Facility.List(page(c), "")
	if err != nil {
		response.InternalError(c, "query failed")
		return
	}
	response.PageOK(c, items, total, page(c).Page, page(c).PageSize)
}

func (h *Handler) ListFacilitiesByOrg(c *gin.Context) {
	id, ok := idParam(c)
	if !ok {
		return
	}
	p := page(c)
	items, total, err := h.Facility.ListByOrganization(p, id)
	if err != nil {
		response.InternalError(c, "query failed")
		return
	}
	response.PageOK(c, items, total, p.Page, p.PageSize)
}

func (h *Handler) CreateFacility(c *gin.Context) {
	var f model.Facility
	if err := c.ShouldBindJSON(&f); err != nil {
		response.Fail(c, 400, "invalid request: "+err.Error())
		return
	}
	if err := h.Facility.Create(&f); err != nil {
		response.Fail(c, 400, err.Error())
		return
	}
	response.Created(c, &f)
}

func (h *Handler) GetFacility(c *gin.Context) {
	id, ok := idParam(c)
	if !ok {
		return
	}
	f, err := h.Facility.Get(id)
	if err != nil {
		response.NotFound(c, "facility")
		return
	}
	response.OK(c, f)
}

func (h *Handler) UpdateFacility(c *gin.Context) {
	id, ok := idParam(c)
	if !ok {
		return
	}
	var f model.Facility
	if err := c.ShouldBindJSON(&f); err != nil {
		response.Fail(c, 400, "invalid request: "+err.Error())
		return
	}
	f.ID = id
	if err := h.Facility.Update(&f); err != nil {
		response.Fail(c, 400, err.Error())
		return
	}
	response.OK(c, gin.H{"message": "facility updated"})
}

func (h *Handler) DeleteFacility(c *gin.Context) {
	id, ok := idParam(c)
	if !ok {
		return
	}
	if err := h.Facility.Delete(id); err != nil {
		response.InternalError(c, "delete failed")
		return
	}
	response.OK(c, gin.H{"message": "facility deleted"})
}

// ----- Emission Sources -----

func (h *Handler) ListEmissionSources(c *gin.Context) {
	items, total, err := h.Source.List(page(c), "")
	if err != nil {
		response.InternalError(c, "query failed")
		return
	}
	response.PageOK(c, items, total, page(c).Page, page(c).PageSize)
}

func (h *Handler) CreateEmissionSource(c *gin.Context) {
	var s model.EmissionSource
	if err := c.ShouldBindJSON(&s); err != nil {
		response.Fail(c, 400, "invalid request: "+err.Error())
		return
	}
	if err := h.Source.Create(&s); err != nil {
		response.Fail(c, 400, err.Error())
		return
	}
	response.Created(c, &s)
}

func (h *Handler) GetEmissionSource(c *gin.Context) {
	id, ok := idParam(c)
	if !ok {
		return
	}
	s, err := h.Source.Get(id)
	if err != nil {
		response.NotFound(c, "emission source")
		return
	}
	response.OK(c, s)
}

func (h *Handler) UpdateEmissionSource(c *gin.Context) {
	id, ok := idParam(c)
	if !ok {
		return
	}
	var s model.EmissionSource
	if err := c.ShouldBindJSON(&s); err != nil {
		response.Fail(c, 400, "invalid request: "+err.Error())
		return
	}
	s.ID = id
	if err := h.Source.Update(&s); err != nil {
		response.Fail(c, 400, err.Error())
		return
	}
	response.OK(c, gin.H{"message": "emission source updated"})
}

func (h *Handler) DeleteEmissionSource(c *gin.Context) {
	id, ok := idParam(c)
	if !ok {
		return
	}
	if err := h.Source.Delete(id); err != nil {
		response.InternalError(c, "delete failed")
		return
	}
	response.OK(c, gin.H{"message": "emission source deleted"})
}

// ----- Emission Factors -----

func (h *Handler) ListEmissionFactors(c *gin.Context) {
	items, total, err := h.Factor.List(page(c), "")
	if err != nil {
		response.InternalError(c, "query failed")
		return
	}
	response.PageOK(c, items, total, page(c).Page, page(c).PageSize)
}

func (h *Handler) CreateEmissionFactor(c *gin.Context) {
	var f model.EmissionFactor
	if err := c.ShouldBindJSON(&f); err != nil {
		response.Fail(c, 400, "invalid request: "+err.Error())
		return
	}
	if err := h.Factor.Create(&f); err != nil {
		response.Fail(c, 400, err.Error())
		return
	}
	response.Created(c, &f)
}

func (h *Handler) GetEmissionFactor(c *gin.Context) {
	id, ok := idParam(c)
	if !ok {
		return
	}
	f, err := h.Factor.Get(id)
	if err != nil {
		response.NotFound(c, "emission factor")
		return
	}
	response.OK(c, f)
}

func (h *Handler) UpdateEmissionFactor(c *gin.Context) {
	id, ok := idParam(c)
	if !ok {
		return
	}
	var f model.EmissionFactor
	if err := c.ShouldBindJSON(&f); err != nil {
		response.Fail(c, 400, "invalid request: "+err.Error())
		return
	}
	f.ID = id
	if err := h.Factor.Update(&f); err != nil {
		response.Fail(c, 400, err.Error())
		return
	}
	response.OK(c, gin.H{"message": "emission factor updated"})
}

func (h *Handler) DeleteEmissionFactor(c *gin.Context) {
	id, ok := idParam(c)
	if !ok {
		return
	}
	if err := h.Factor.Delete(id); err != nil {
		response.InternalError(c, "delete failed")
		return
	}
	response.OK(c, gin.H{"message": "emission factor deleted"})
}
