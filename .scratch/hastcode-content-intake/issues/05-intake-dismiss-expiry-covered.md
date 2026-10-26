# 05: Intake: dismiss, expiry, covered Leads and the run summary

**What to build:** The queue only holds open decisions. `make lead-dismiss ID=…` hides a Lead for good. Open Leads first seen more than 7 days before the run's clock expire automatically. A Lead whose normalised link already appears in any Piece's Sources or body is marked covered in the queue. The run summary shows new Leads per Vertical and per suggested kind, alongside the failures and silent Sources from ticket 03.

Spec: `../spec.md` ("Intake", stories 17–19, 24).

**Blocked by:** 02 (Intake tracer)

**Status:** resolved

- [x] Intake-seam tests: dismissed Lead absent from later queues; Lead expires at 7 days and not at 6; covered marking from a fixture content tree; summary counts per Vertical and kind
- [x] Seen-state records each Lead's decision (open, picked, dismissed, expired) and first-seen date
- [x] README documents dismissing and expiry
