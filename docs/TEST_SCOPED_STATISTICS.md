# Manual Test Plan for Scoped Statistics Bug Fix

## Bug Description
When a user with `books_read_by_my_role` permission views the books page, the data table correctly shows filtered results (e.g., "No results found" if no books match the scope), but the statistics cards and filter badges still show counts for ALL books in the system.

## Root Cause
The `GetCountByFilter` method in `GenericPageController` was using `GetListAdvanced` which does not apply scope filtering, instead of `GetList` which does.

## Fix Applied
Changed line 263 in `/app/contracts/generic_page_controller.go`:
```go
// OLD: result, err := c.service.GetListAdvanced(req, filters)
// NEW: result, err := c.service.GetList(req)
```

## Manual Testing Steps

### Test Case 1: Member with by_my_role Permission
1. Login as `lucky@test.com` (Member role)
2. Ensure Member role has `books_read_by_my_role` permission
3. Navigate to `/admin/books`
4. **Expected Results:**
   - Total Books card should show 0 (not total system books)
   - Available filter badge should show 0 (not 52)
   - Borrowed filter badge should show 0 (not 30)
   - Maintenance filter badge should show 0 (not 16)
   - Data table should show "No results found"

### Test Case 2: Editor with by_my_role Permission
1. Create some books as an editor user
2. Login as another editor user
3. Navigate to `/admin/books`
4. **Expected Results:**
   - Statistics should only count books created by users with Editor role
   - Filter badges should match the actual filtered counts
   - Data table should only show books from editors

### Test Case 3: Admin with by_all Permission
1. Login as admin user
2. Navigate to `/admin/books`
3. **Expected Results:**
   - Statistics should show ALL books in the system
   - Filter badges should show total counts for each status
   - Data table should show all books

### Test Case 4: Creating Books and Verifying Statistics Update
1. Login as a member with `books_read_by_my_role`
2. Create a new book
3. Have another member create a book
4. **Expected Results:**
   - Statistics should now show 2 total books
   - Filter badges should update to reflect the new books
   - Both members should see both books (same role)

## Regression Test
A comprehensive regression test has been created at:
`/tests/feature/scoped_statistics_regression_test.go`

This test covers:
- Admin seeing all statistics
- Editor seeing only editor statistics
- Member seeing only member statistics
- Filter badges matching actual data counts
- Users with no books in scope seeing zero statistics

## Additional Considerations
- The fix ensures that ALL statistics queries respect scope filtering
- This applies to any service that has scope filtering enabled
- Currently only BookService has scope filtering enabled via `.EnableScopeFiltering(auth.ServiceBooks, "created_by")`