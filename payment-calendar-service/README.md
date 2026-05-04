# Payment Calendar Service

Payment Calendar Service stores planned project incomes and expenses, marks actual payments as paid, computes overdue payments from planned dates, and returns project financial summaries.

## Domain rules

- Each payment belongs to a project and may belong to a project stage.
- Payment types are `income` and `expense`.
- Stored statuses are `planned`, `paid`, and `cancelled`.
- Overdue is not stored in the database. A payment is overdue when it is still `planned` and its `planned_date` is before the current date.
- Project and stage references are validated through `project-service` over gRPC.

## gRPC methods

- `CreatePayment`
- `GetPayment`
- `ListPayments`
- `UpdatePayment`
- `DeletePayment`
- `MarkPaymentPaid`
- `GetProjectSummary`
