# 04: Prune command for the Image Store

**What to build:** The Author runs one command that lists every object in the R2 bucket that the current content would not produce, prints them with sizes and a total, and deletes nothing unless `--delete` is passed. It reads the same config and environment as the builder, computes the set of keys a production build of the current content tree would upload (without uploading), lists the bucket, and diffs. Deletes go in batches.

Spec: `../spec.md` ("Image Store"). Research: `docs/research/cloudflare-r2-and-deploy.md` section 4.

**Blocked by:** 03 (Image Store: R2 mode)

**Status:** ready-for-agent

- [ ] Tests against the fake store: with extra objects present, dry run lists exactly the orphans and deletes nothing; `--delete` removes exactly them and leaves referenced keys
- [ ] Objects a Preview could have uploaded for a draft are listed as orphans only when that draft is not in the current tree (document that Preview images may be pruned)
- [ ] Missing credentials fail with the same message as the builder
- [ ] Documented in README with the warning that it is manual and dry-run by default
