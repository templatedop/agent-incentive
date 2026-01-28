-- ============================================================================
-- Migration Rollback: 002_add_commission_batches
-- Description: Remove commission_batches table
-- ============================================================================

DROP TABLE IF EXISTS commission_batches CASCADE;
