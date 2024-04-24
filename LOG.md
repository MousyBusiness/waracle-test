# Infrastructure Decision Log

### Purpose

This document gives a brief record of the __assumptions__ made behind a technology choice. Hopefully mistakes aren't made twice and improvements can be made if
the underlying assumptions are proven invalid or outdated.

---

#### [24-04-2024] CloudRun for compute instead of Cloud Functions

CloudRun provides a superior development experience as well as facilitates an offline testing setup using Docker. CloudRun costs and execution speeds are on par
with Cloud Functions.

#### [24-04-2024] API Key authentication for APIs

Given the time constraints of the waracle-test, API keys will be used for authentication instead of OAuth2. This provides an additional advantage of
demonstration of security concepts when working with secrets in GCP.

