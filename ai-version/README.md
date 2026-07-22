# AI vs Me

## Similarities

- Both implementations use Go's standard net/http package.
- Both organize the code into packages.
- Both implement CRUD endpoints.
- Both store tasks in memory using a slice.
- Both return JSON responses.

---

## Differences

### My implementation

- Added Swagger documentation.
- Better validation.
- Better route organization.
- More descriptive error handling.
- More production-ready.

### AI implementation

- Simpler code.
- Shorter handlers.
- Less validation.
- Simpler error responses.
- Easier to read but less robust.

---

## What AI did well

- Quickly generated a working CRUD API.
- Produced clean and readable code.
- Correctly organized the project into packages.

---

## What AI got wrong

- Initially assumed the wrong Go module import path.
- Required manual review and testing.
- Did not include the same level of validation and documentation as my implementation.

---

## Lessons Learned

AI significantly speeds up development, but generated code still requires careful review, testing, debugging, and improvement before it is production-ready.