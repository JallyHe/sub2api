package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/dialect"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// CreditLedger is an append-only log of credit transactions per user.
// Never update or delete rows; the balance_after snapshot allows auditing.
type CreditLedger struct {
	ent.Schema
}

func (CreditLedger) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "credit_ledger"},
	}
}

func (CreditLedger) Fields() []ent.Field {
	return []ent.Field{
		field.Int64("user_id").
			Comment("关联用户 ID"),
		field.Int64("delta").
			Comment("积分变动量，正数为充值，负数为消耗"),
		field.String("reason").
			MaxLen(50).
			NotEmpty().
			Comment("变动原因：purchase | api_call | admin_grant | expiry_reset"),
		field.String("ref_id").
			MaxLen(100).
			Optional().
			Nillable().
			Comment("关联记录 ID，如订单 ID 或请求 ID"),
		field.Int64("balance_after").
			Comment("操作后的余额快照"),
		field.String("model").
			MaxLen(100).
			Optional().
			Nillable().
			Comment("消耗时的模型名称，reason=api_call 时填写"),
		field.Time("created_at").
			Immutable().
			Default(time.Now).
			SchemaType(map[string]string{dialect.Postgres: "timestamptz"}),
	}
}

func (CreditLedger) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("user_id", "created_at"),
		index.Fields("user_id"),
	}
}
