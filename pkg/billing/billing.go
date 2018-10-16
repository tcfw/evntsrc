package billing

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"github.com/tcfw/evntsrc/pkg/users/protos"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	stripe "github.com/stripe/stripe-go"
	"github.com/stripe/stripe-go/customer"
	"github.com/stripe/stripe-go/plan"
	"github.com/stripe/stripe-go/product"
	"github.com/stripe/stripe-go/sub"
	pb "github.com/tcfw/evntsrc/pkg/billing/protos"
)

const UserMetadataCustomerId = "stripe_customer_id"

//Server core struct
type Server struct {
	mu sync.Mutex
}

//NewServer creates a new struct to interface the streams server
func NewServer() *Server {
	return &Server{}
}

//GetPlans @TODO
func (s *Server) GetPlans(ctx context.Context, _ *pb.Empty) (*pb.PlanList, error) {

	plans := []*pb.Plan{}

	params := &stripe.PlanListParams{Active: stripe.Bool(true)}
	i := plan.List(params)
	for i.Next() {
		p := i.Plan()

		product, err := product.Get(p.Product.ID, nil)
		if err != nil {
			return nil, err
		}

		plans = append(plans, &pb.Plan{
			Active:         p.Active,
			AggregateUsage: p.AggregateUsage,
			Amount:         p.Amount,
			BillingScheme:  string(p.BillingScheme),
			Currency:       string(p.Currency),
			Deleted:        p.Deleted,
			Id:             p.ID,
			Interval:       string(p.Interval),
			IntervalCount:  p.IntervalCount,
			LiveMode:       p.Livemode,
			Metadata:       p.Metadata,
			Nickname:       p.Nickname,
			Product: &pb.Product{
				Active:      product.Active,
				Attributes:  product.Attributes,
				Caption:     product.Caption,
				Description: product.Description,
				Id:          product.ID,
				Metadata:    product.Metadata,
				Name:        product.Name,
			},
			TiresMode:       p.TiersMode,
			TrialPeriodDays: p.TrialPeriodDays,
			UsageType:       string(p.UsageType),
		})
	}

	return &pb.PlanList{Plans: plans}, nil
}

//CreateCustomerFromUser @TODO
func (s *Server) CreateCustomerFromUser(ctx context.Context, request *pb.CreateRequest) (*pb.CreateResponse, error) {

	if request.CardToken == "" || request.UserId == "" {
		return nil, errors.New("All params required")
	}

	users, err := newUserClient(ctx)
	if err != nil {
		return nil, err
	}

	user, err := users.Find(ctx, &evntsrc_users.UserRequest{Query: &evntsrc_users.UserRequest_Id{Id: request.UserId}})
	if err != nil {
		return nil, err
	}

	stripeParams := &stripe.CustomerParams{
		Description: &user.Name,
		Email:       &user.Email,
	}
	if request.CardToken != "" {
		stripeParams.SetSource(request.CardToken)
	}

	customer, err := customer.New(stripeParams)
	if err != nil {
		return nil, err
	}

	user.Metadata[UserMetadataCustomerId] = []byte(customer.ID)
	_, err = users.Update(ctx, &evntsrc_users.UserUpdateRequest{Id: user.Id, User: user})
	if err != nil {
		return nil, fmt.Errorf("Failed to set customer id to user (%s => %s)", user.Id, customer.ID)
	}

	return &pb.CreateResponse{CustomerId: customer.ID}, nil
}

//AttachPaymentMethod @TODO
func (s *Server) AttachPaymentMethod(ctx context.Context, request *pb.CreateRequest) (*pb.Empty, error) {
	if request.CardToken == "" || request.UserId == "" {
		return nil, errors.New("All params required")
	}

	users, err := newUserClient(ctx)
	if err != nil {
		return nil, err
	}

	user, err := users.Find(ctx, &evntsrc_users.UserRequest{Query: &evntsrc_users.UserRequest_Id{Id: request.UserId}})
	if err != nil {
		return nil, err
	}

	customerID, ok := user.Metadata[UserMetadataCustomerId]
	if !ok {
		_, err := s.CreateCustomerFromUser(ctx, request)
		return nil, err
	}
	params := &stripe.CustomerParams{}
	params.SetSource(request.CardToken)

	_, err = customer.Update(string(customerID), params)
	return nil, err
}

//GetUserInfo @TODO
func (s *Server) GetUserInfo(ctx context.Context, request *pb.InfoRequest) (*pb.Customer, error) {

	if request.UserId == "" {
		return nil, errors.New("User id required")
	}

	users, err := newUserClient(ctx)
	if err != nil {
		return nil, err
	}

	user, err := users.Find(ctx, &evntsrc_users.UserRequest{Query: &evntsrc_users.UserRequest_Id{Id: request.UserId}})
	if err != nil {
		return nil, err
	}

	customerID, ok := user.Metadata[UserMetadataCustomerId]
	if !ok {
		createResp, err := s.CreateCustomerFromUser(ctx, &pb.CreateRequest{UserId: request.UserId})
		if err != nil {
			return nil, err
		}
		customerID = []byte(createResp.CustomerId)
	}

	customer, err := customer.Get(string(customerID), nil)
	if err != nil {
		return nil, err
	}

	rCustomer := &pb.Customer{
		AccountBalance: customer.AccountBalance,
		Created:        customer.Created,
		Currency:       string(customer.Currency),
		Deleted:        customer.Deleted,
		Delinquent:     customer.Delinquent,
		Description:    customer.Description,
		Email:          customer.Email,
		Id:             customer.ID,
		InvoicePrefix:  customer.InvoicePrefix,
		LiveMode:       customer.Livemode,
		Metadata:       customer.Metadata,
		Subscriptions:  []*pb.Subscription{},
	}

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		subs := sub.List(&stripe.SubscriptionListParams{Customer: string(customerID)})
		for subs.Next() {
			subscription := subs.Subscription()

			sub := &pb.Subscription{
				Billing:             string(subscription.Billing),
				BillingCycleAnchor:  subscription.BillingCycleAnchor,
				CanceledAt:          subscription.CanceledAt,
				CurrentPeriodEnds:   subscription.CurrentPeriodEnd,
				CurrentPeriodStarts: subscription.CurrentPeriodStart,
				DaysUntilDue:        subscription.DaysUntilDue,
				CancelAtPeriodEnd:   subscription.CancelAtPeriodEnd,
				EndedAt:             subscription.EndedAt,
				Id:                  subscription.ID,
				Metadata:            subscription.Metadata,
				Plan:                subscription.Plan.ID,
				Start:               subscription.Start,
				Status:              string(subscription.Status),
				TaxPercent:          subscription.TaxPercent,
				TrialEnds:           subscription.TrialEnd,
				TrialStarts:         subscription.TrialStart,
			}

			if subscription.Discount != nil && subscription.Discount.Coupon != nil {
				if subscription.Discount.Coupon.AmountOff > 0 {
					sub.Discount = fmt.Sprintf("$%.2f", float64(subscription.Discount.Coupon.AmountOff)/100)
				} else if subscription.Discount.Coupon.PercentOff > 0 {
					sub.Discount = "%" + fmt.Sprintf("%.0f", subscription.Discount.Coupon.PercentOff) + " off"
				} else {
					sub.Discount = subscription.Discount.Coupon.Name
				}
				sub.DiscountEnds = subscription.Discount.End
			}

			rCustomer.Subscriptions = append(rCustomer.Subscriptions, sub)
		}
		wg.Done()
	}()

	go func() {
		rCustomer.Sources = []*pb.Source{}
		for _, stripeSource := range customer.Sources.Data {
			source := &pb.Source{
				Id:   stripeSource.ID,
				Type: string(stripeSource.Type),
			}

			switch stripeSource.Type {
			case "card":
				source.Source = &pb.Source_Card{
					Card: &pb.Card{
						Id:       stripeSource.Card.ID,
						Brand:    string(stripeSource.Card.Brand),
						Country:  stripeSource.Card.Country,
						ExpMonth: uint32(stripeSource.Card.ExpMonth),
						ExpYear:  uint32(stripeSource.Card.ExpYear),
						Name:     stripeSource.Card.Name,
						CvcCheck: string(stripeSource.Card.CVCCheck),
						Number:   stripeSource.Card.Last4,
					},
				}
			}

			rCustomer.Sources = append(rCustomer.Sources, source)
		}
		wg.Done()
	}()

	wg.Wait()

	return rCustomer, nil
}

//GetUserSubscriptions @TODO
func (s *Server) GetUserSubscriptions(ctx context.Context, request *pb.UserSubscriptionsRequest) (*pb.UserSubscriptionsResponse, error) {
	return nil, status.Errorf(codes.Unavailable, "Not implemented")
}

//SubscribeUser @TODO
func (s *Server) SubscribeUser(ctx context.Context, request *pb.SubscribeRequest) (*pb.SubscribeResponse, error) {
	return nil, status.Errorf(codes.Unavailable, "Not implemented")
}

//TerminateSubscription @TODO
func (s *Server) TerminateSubscription(ctx context.Context, request *pb.TerminateRequest) (*pb.Empty, error) {
	return nil, status.Errorf(codes.Unavailable, "Not implemented")
}
