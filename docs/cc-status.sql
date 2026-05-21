CREATE TABLE `usage_reports` (`id` integer PRIMARY KEY AUTOINCREMENT,`client_id` text NOT NULL,`request_id` text NOT NULL,`app_type` text NOT NULL,`model` text NOT NULL,`input_tokens` integer NOT NULL DEFAULT 0,`output_tokens` integer NOT NULL DEFAULT 0,`cache_read_tokens` integer NOT NULL DEFAULT 0,`cache_creation_tokens` integer NOT NULL DEFAULT 0,`input_cost_usd` decimal(20,10) NOT NULL DEFAULT "0",`output_cost_usd` decimal(20,10) NOT NULL DEFAULT "0",`cache_read_cost_usd` decimal(20,10) NOT NULL DEFAULT "0",`cache_creation_cost_usd` decimal(20,10) NOT NULL DEFAULT "0",`total_cost_usd` decimal(20,10) NOT NULL DEFAULT "0",`session_id` text,`pricing_source` text NOT NULL,`created_at` integer NOT NULL,`data_source` text NOT NULL,`inserted_at` datetime NOT NULL);

CREATE INDEX `idx_usage_reports_client_id` ON `usage_reports`(`client_id`);

CREATE UNIQUE INDEX `idx_usage_reports_client_request` ON `usage_reports`(`client_id`,`request_id`);

CREATE INDEX `idx_usage_reports_created_at` ON `usage_reports`(`created_at`);

CREATE INDEX `idx_usage_reports_model` ON `usage_reports`(`model`);


CREATE TABLE `model_pricing` (`id` integer PRIMARY KEY AUTOINCREMENT,`model_id` text NOT NULL,`display_name` text,`input_cost_per_million` decimal(20,10) NOT NULL,`output_cost_per_million` decimal(20,10) NOT NULL,`cache_read_cost_per_million` decimal(20,10) NOT NULL DEFAULT "0",`cache_creation_cost_per_million` decimal(20,10) NOT NULL DEFAULT "0",`is_placeholder` numeric NOT NULL DEFAULT false,`created_at` datetime NOT NULL,`updated_at` datetime NOT NULL);

CREATE INDEX `idx_model_pricing_is_placeholder` ON `model_pricing`(`is_placeholder`);

CREATE UNIQUE INDEX `idx_model_pricing_model_id` ON `model_pricing`(`model_id`);