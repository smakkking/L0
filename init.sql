CREATE TABLE WB_Order (
  "order_uid" text,
  "track_number" text,
  "entry" text,

  "delivery.name" text,
  "delivery.phone" text,
  "delivery.zip" text,
  "delivery.city" text,
  "delivery.address" text,
  "delivery.region" text,
  "delivery.email" text,
  
  "payment.transaction" text,
  "payment.request_id" text,
  "payment.currency" text,
  "payment.provider" text,
  "payment.amount" bigint,
  "payment.payment_dt" bigint,
  "payment.bank" text,
  "payment.delivery_cost" bigint,
  "payment.goods_total" bigint,
  "payment.custom_fee" bigint,
  
  "items" jsonb,
  
  "locale" text,
  "internal_signature" text,
  "customer_id" text,
  "delivery_service" text,
  "shardkey" text,
  "sm_id" bigint,
  "date_created" text,
  "oof_shard" text
);