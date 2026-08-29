package services

import (
	"context"
	"github.com/ngwebwork/payflow/internal/repositories"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type AnalyticsService struct {
	r *repositories.AnalyticsRepository
}

func NewAnalyticsService(r *repositories.AnalyticsRepository) *AnalyticsService {
	return &AnalyticsService{r}
}
func period(v string) time.Time {
	n := time.Now().UTC()
	switch v {
	case "7d":
		return n.AddDate(0, 0, -7)
	case "3m":
		return n.AddDate(0, -3, 0)
	case "6m":
		return n.AddDate(0, -6, 0)
	case "1y":
		return n.AddDate(-1, 0, 0)
	}
	return n.AddDate(0, 0, -30)
}
func (s *AnalyticsService) Revenue(c context.Context, u primitive.ObjectID, p string) ([]bson.M, error) {
	return s.r.Series(c, "transactions", u, period(p))
}
func (s *AnalyticsService) Transactions(c context.Context, u primitive.ObjectID, p string) (bson.M, error) {
	f := bson.M{"createdAt": bson.M{"$gte": period(p)}}
	out := bson.M{}
	for _, v := range []string{"successful", "pending", "failed", "refunded"} {
		n, e := s.r.Count(c, "transactions", u, bson.M{"status": v, "createdAt": f["createdAt"]})
		if e != nil {
			return nil, e
		}
		out[v] = n
	}
	out["total"], _ = s.r.Count(c, "transactions", u, f)
	return out, nil
}
func (s *AnalyticsService) Customers(c context.Context, u primitive.ObjectID, p string) (bson.M, error) {
	t, e := s.r.Count(c, "customers", u, bson.M{})
	if e != nil {
		return nil, e
	}
	n, e := s.r.Count(c, "customers", u, bson.M{"createdAt": bson.M{"$gte": period(p)}})
	return bson.M{"totalCustomers": t, "newCustomers": n, "customerGrowth": n}, e
}
func (s *AnalyticsService) Payments(c context.Context, u primitive.ObjectID, p string) (bson.M, error) {
	x, e := s.Transactions(c, u, p)
	if e != nil {
		return nil, e
	}
	for _, v := range []string{"successful", "pending", "failed", "refunded"} {
		n, er := s.r.Count(c, "payments", u, bson.M{"status": v, "createdAt": bson.M{"$gte": period(p)}})
		if er != nil {
			return nil, er
		}
		x[v] = n
	}
	total := x["successful"].(int64) + x["pending"].(int64) + x["failed"].(int64) + x["refunded"].(int64)
	x["total"] = total
	if total > 0 {
		x["paymentSuccessRate"] = float64(x["successful"].(int64)) * 100 / float64(total)
	} else {
		x["paymentSuccessRate"] = 0
	}
	return x, nil
}
