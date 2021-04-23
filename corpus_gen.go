// This file was automatically generated by genny.
// Any changes will be lost if this file is regenerated.
// see https://github.com/cheekybits/genny

package fball

import (
	"context"
	"encoding/json"
	"time"
)

func (c Corpus) getTimezoneResponse(
	ctx context.Context, endpoint string, max int, rng tRange, policy refreshPolicy,
	params urlQueryStringer) ([]TimezoneResponse, error) {
	// Query the countries from the database.
	resp := []TimezoneResponse{}

	q1 := time.Now()
	err := c.cache.Query(ctx, endpoint, params, max, rng, func(data []byte) error {
		cr := TimezoneResponse{}
		if err := json.Unmarshal(data, &cr); err != nil {
			return err
		}
		resp = append(resp, cr)
		return nil
	})
	q2 := time.Now()
	c.logger.Printf("INFO - %q query time: %dms", endpoint, q2.Sub(q1)/time.Millisecond)

	if err == nil && len(resp) != 0 && policy.Valid(time.Now(), resp[0].When()) {
		return resp, nil
	} else if err != nil {
		c.logger.Printf("WARNING - query error for countries: %v", err)
	}

	// Either the data is not available or it has expired.
	rQ := TimezoneResponse{}

	s1 := time.Now()
	err = c.fballc.Get(ctx, endpoint, &rQ, params)
	s2 := time.Now()
	c.logger.Printf("INFO - %q api call time: %dms", endpoint, s2.Sub(s1)/time.Millisecond)

	if err != nil {
		// We tolerate stale data if the api call fails. We still log it.
		if len(resp) != 0 {
			c.logger.Printf("WARNING - unable to query countries: %v.", err)
			c.logger.Printf("WARNING - returning stale data for countries.")
			return resp, nil
		} else {
			return nil, err
		}
	}

	i1 := time.Now()
	err = c.cache.Insert(ctx, endpoint, &rQ, params)
	i2 := time.Now()
	c.logger.Printf("INFO - %q insert time: %dms", endpoint, i2.Sub(i1)/time.Millisecond)

	if err != nil {
		c.logger.Printf("ERROR - unable to write country to cache: %v", err)
	}

	return []TimezoneResponse{rQ}, nil
}

func (c Corpus) getCountryResponse(
	ctx context.Context, endpoint string, max int, rng tRange, policy refreshPolicy,
	params urlQueryStringer) ([]CountryResponse, error) {
	// Query the countries from the database.
	resp := []CountryResponse{}

	q1 := time.Now()
	err := c.cache.Query(ctx, endpoint, params, max, rng, func(data []byte) error {
		cr := CountryResponse{}
		if err := json.Unmarshal(data, &cr); err != nil {
			return err
		}
		resp = append(resp, cr)
		return nil
	})
	q2 := time.Now()
	c.logger.Printf("INFO - %q query time: %dms", endpoint, q2.Sub(q1)/time.Millisecond)

	if err == nil && len(resp) != 0 && policy.Valid(time.Now(), resp[0].When()) {
		return resp, nil
	} else if err != nil {
		c.logger.Printf("WARNING - query error for countries: %v", err)
	}

	// Either the data is not available or it has expired.
	rQ := CountryResponse{}

	s1 := time.Now()
	err = c.fballc.Get(ctx, endpoint, &rQ, params)
	s2 := time.Now()
	c.logger.Printf("INFO - %q api call time: %dms", endpoint, s2.Sub(s1)/time.Millisecond)

	if err != nil {
		// We tolerate stale data if the api call fails. We still log it.
		if len(resp) != 0 {
			c.logger.Printf("WARNING - unable to query countries: %v.", err)
			c.logger.Printf("WARNING - returning stale data for countries.")
			return resp, nil
		} else {
			return nil, err
		}
	}

	i1 := time.Now()
	err = c.cache.Insert(ctx, endpoint, &rQ, params)
	i2 := time.Now()
	c.logger.Printf("INFO - %q insert time: %dms", endpoint, i2.Sub(i1)/time.Millisecond)

	if err != nil {
		c.logger.Printf("ERROR - unable to write country to cache: %v", err)
	}

	return []CountryResponse{rQ}, nil
}

func (c Corpus) getSeasonResponse(
	ctx context.Context, endpoint string, max int, rng tRange, policy refreshPolicy,
	params urlQueryStringer) ([]SeasonResponse, error) {
	// Query the countries from the database.
	resp := []SeasonResponse{}

	q1 := time.Now()
	err := c.cache.Query(ctx, endpoint, params, max, rng, func(data []byte) error {
		cr := SeasonResponse{}
		if err := json.Unmarshal(data, &cr); err != nil {
			return err
		}
		resp = append(resp, cr)
		return nil
	})
	q2 := time.Now()
	c.logger.Printf("INFO - %q query time: %dms", endpoint, q2.Sub(q1)/time.Millisecond)

	if err == nil && len(resp) != 0 && policy.Valid(time.Now(), resp[0].When()) {
		return resp, nil
	} else if err != nil {
		c.logger.Printf("WARNING - query error for countries: %v", err)
	}

	// Either the data is not available or it has expired.
	rQ := SeasonResponse{}

	s1 := time.Now()
	err = c.fballc.Get(ctx, endpoint, &rQ, params)
	s2 := time.Now()
	c.logger.Printf("INFO - %q api call time: %dms", endpoint, s2.Sub(s1)/time.Millisecond)

	if err != nil {
		// We tolerate stale data if the api call fails. We still log it.
		if len(resp) != 0 {
			c.logger.Printf("WARNING - unable to query countries: %v.", err)
			c.logger.Printf("WARNING - returning stale data for countries.")
			return resp, nil
		} else {
			return nil, err
		}
	}

	i1 := time.Now()
	err = c.cache.Insert(ctx, endpoint, &rQ, params)
	i2 := time.Now()
	c.logger.Printf("INFO - %q insert time: %dms", endpoint, i2.Sub(i1)/time.Millisecond)

	if err != nil {
		c.logger.Printf("ERROR - unable to write country to cache: %v", err)
	}

	return []SeasonResponse{rQ}, nil
}

func (c Corpus) getLeagueInfoResponse(
	ctx context.Context, endpoint string, max int, rng tRange, policy refreshPolicy,
	params urlQueryStringer) ([]LeagueInfoResponse, error) {
	// Query the countries from the database.
	resp := []LeagueInfoResponse{}

	q1 := time.Now()
	err := c.cache.Query(ctx, endpoint, params, max, rng, func(data []byte) error {
		cr := LeagueInfoResponse{}
		if err := json.Unmarshal(data, &cr); err != nil {
			return err
		}
		resp = append(resp, cr)
		return nil
	})
	q2 := time.Now()
	c.logger.Printf("INFO - %q query time: %dms", endpoint, q2.Sub(q1)/time.Millisecond)

	if err == nil && len(resp) != 0 && policy.Valid(time.Now(), resp[0].When()) {
		return resp, nil
	} else if err != nil {
		c.logger.Printf("WARNING - query error for countries: %v", err)
	}

	// Either the data is not available or it has expired.
	rQ := LeagueInfoResponse{}

	s1 := time.Now()
	err = c.fballc.Get(ctx, endpoint, &rQ, params)
	s2 := time.Now()
	c.logger.Printf("INFO - %q api call time: %dms", endpoint, s2.Sub(s1)/time.Millisecond)

	if err != nil {
		// We tolerate stale data if the api call fails. We still log it.
		if len(resp) != 0 {
			c.logger.Printf("WARNING - unable to query countries: %v.", err)
			c.logger.Printf("WARNING - returning stale data for countries.")
			return resp, nil
		} else {
			return nil, err
		}
	}

	i1 := time.Now()
	err = c.cache.Insert(ctx, endpoint, &rQ, params)
	i2 := time.Now()
	c.logger.Printf("INFO - %q insert time: %dms", endpoint, i2.Sub(i1)/time.Millisecond)

	if err != nil {
		c.logger.Printf("ERROR - unable to write country to cache: %v", err)
	}

	return []LeagueInfoResponse{rQ}, nil
}

func (c Corpus) getTeamInfoResponse(
	ctx context.Context, endpoint string, max int, rng tRange, policy refreshPolicy,
	params urlQueryStringer) ([]TeamInfoResponse, error) {
	// Query the countries from the database.
	resp := []TeamInfoResponse{}

	q1 := time.Now()
	err := c.cache.Query(ctx, endpoint, params, max, rng, func(data []byte) error {
		cr := TeamInfoResponse{}
		if err := json.Unmarshal(data, &cr); err != nil {
			return err
		}
		resp = append(resp, cr)
		return nil
	})
	q2 := time.Now()
	c.logger.Printf("INFO - %q query time: %dms", endpoint, q2.Sub(q1)/time.Millisecond)

	if err == nil && len(resp) != 0 && policy.Valid(time.Now(), resp[0].When()) {
		return resp, nil
	} else if err != nil {
		c.logger.Printf("WARNING - query error for countries: %v", err)
	}

	// Either the data is not available or it has expired.
	rQ := TeamInfoResponse{}

	s1 := time.Now()
	err = c.fballc.Get(ctx, endpoint, &rQ, params)
	s2 := time.Now()
	c.logger.Printf("INFO - %q api call time: %dms", endpoint, s2.Sub(s1)/time.Millisecond)

	if err != nil {
		// We tolerate stale data if the api call fails. We still log it.
		if len(resp) != 0 {
			c.logger.Printf("WARNING - unable to query countries: %v.", err)
			c.logger.Printf("WARNING - returning stale data for countries.")
			return resp, nil
		} else {
			return nil, err
		}
	}

	i1 := time.Now()
	err = c.cache.Insert(ctx, endpoint, &rQ, params)
	i2 := time.Now()
	c.logger.Printf("INFO - %q insert time: %dms", endpoint, i2.Sub(i1)/time.Millisecond)

	if err != nil {
		c.logger.Printf("ERROR - unable to write country to cache: %v", err)
	}

	return []TeamInfoResponse{rQ}, nil
}

func (c Corpus) getTeamStatsResponse(
	ctx context.Context, endpoint string, max int, rng tRange, policy refreshPolicy,
	params urlQueryStringer) ([]TeamStatsResponse, error) {
	// Query the countries from the database.
	resp := []TeamStatsResponse{}

	q1 := time.Now()
	err := c.cache.Query(ctx, endpoint, params, max, rng, func(data []byte) error {
		cr := TeamStatsResponse{}
		if err := json.Unmarshal(data, &cr); err != nil {
			return err
		}
		resp = append(resp, cr)
		return nil
	})
	q2 := time.Now()
	c.logger.Printf("INFO - %q query time: %dms", endpoint, q2.Sub(q1)/time.Millisecond)

	if err == nil && len(resp) != 0 && policy.Valid(time.Now(), resp[0].When()) {
		return resp, nil
	} else if err != nil {
		c.logger.Printf("WARNING - query error for countries: %v", err)
	}

	// Either the data is not available or it has expired.
	rQ := TeamStatsResponse{}

	s1 := time.Now()
	err = c.fballc.Get(ctx, endpoint, &rQ, params)
	s2 := time.Now()
	c.logger.Printf("INFO - %q api call time: %dms", endpoint, s2.Sub(s1)/time.Millisecond)

	if err != nil {
		// We tolerate stale data if the api call fails. We still log it.
		if len(resp) != 0 {
			c.logger.Printf("WARNING - unable to query countries: %v.", err)
			c.logger.Printf("WARNING - returning stale data for countries.")
			return resp, nil
		} else {
			return nil, err
		}
	}

	i1 := time.Now()
	err = c.cache.Insert(ctx, endpoint, &rQ, params)
	i2 := time.Now()
	c.logger.Printf("INFO - %q insert time: %dms", endpoint, i2.Sub(i1)/time.Millisecond)

	if err != nil {
		c.logger.Printf("ERROR - unable to write country to cache: %v", err)
	}

	return []TeamStatsResponse{rQ}, nil
}

func (c Corpus) getVenueResponse(
	ctx context.Context, endpoint string, max int, rng tRange, policy refreshPolicy,
	params urlQueryStringer) ([]VenueResponse, error) {
	// Query the countries from the database.
	resp := []VenueResponse{}

	q1 := time.Now()
	err := c.cache.Query(ctx, endpoint, params, max, rng, func(data []byte) error {
		cr := VenueResponse{}
		if err := json.Unmarshal(data, &cr); err != nil {
			return err
		}
		resp = append(resp, cr)
		return nil
	})
	q2 := time.Now()
	c.logger.Printf("INFO - %q query time: %dms", endpoint, q2.Sub(q1)/time.Millisecond)

	if err == nil && len(resp) != 0 && policy.Valid(time.Now(), resp[0].When()) {
		return resp, nil
	} else if err != nil {
		c.logger.Printf("WARNING - query error for countries: %v", err)
	}

	// Either the data is not available or it has expired.
	rQ := VenueResponse{}

	s1 := time.Now()
	err = c.fballc.Get(ctx, endpoint, &rQ, params)
	s2 := time.Now()
	c.logger.Printf("INFO - %q api call time: %dms", endpoint, s2.Sub(s1)/time.Millisecond)

	if err != nil {
		// We tolerate stale data if the api call fails. We still log it.
		if len(resp) != 0 {
			c.logger.Printf("WARNING - unable to query countries: %v.", err)
			c.logger.Printf("WARNING - returning stale data for countries.")
			return resp, nil
		} else {
			return nil, err
		}
	}

	i1 := time.Now()
	err = c.cache.Insert(ctx, endpoint, &rQ, params)
	i2 := time.Now()
	c.logger.Printf("INFO - %q insert time: %dms", endpoint, i2.Sub(i1)/time.Millisecond)

	if err != nil {
		c.logger.Printf("ERROR - unable to write country to cache: %v", err)
	}

	return []VenueResponse{rQ}, nil
}
