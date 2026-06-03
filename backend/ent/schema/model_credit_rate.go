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

// ModelCreditRate holds per-model credit consumption ratios.
// Admins configure this table; model_pattern supports glob wildcards (e.g. "claude-opus-4*").
// When multiple rows match a model name, the row with the highest priority wins.
type ModelCreditRate struct {
	ent.Schema
}

func (ModelCreditRate) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "model_credit_rates"},
	}
}

func (ModelCreditRate) Fields() []ent.Field {
	return []ent.Field{
		field.String("model_pattern").
			MaxLen(200).
			NotEmpty().
			Comment("模型匹配模式，支持 * 通配符，如 claude-opus-4*"),
		field.Int64("credits_per_1k_tokens_input").
			Default(1).
			Comment("每 1000 个输入 token 消耗的积分数"),
		field.Int64("credits_per_1k_tokens_output").
			Default(3).
			Comment("每 1000 个输出 token 消耗的积分数"),
		field.Int("priority").
			Default(0).
			Comment("匹配优先级，值越大优先级越高"),
		field.Time("created_at").
			Immutable().
			Default(time.Now).
			SchemaType(map[string]string{dialect.Postgres: "timestamptz"}),
		field.Time("updated_at").
			Default(time.Now).
			UpdateDefault(time.Now).
			SchemaType(map[string]string{dialect.Postgres: "timestamptz"}),
	}
}

func (ModelCreditRate) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("priority"),
	}
}
