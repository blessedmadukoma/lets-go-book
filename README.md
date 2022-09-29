#### Let's Go Code Implementation - Alex Edwards
## Chapter 11 - User Authentication:
 - 11.3 User Sign up
   - added bcrypt package for hashing passwords
   - updated insert method to hash password, check if email is duplicate and either return an error if there is any or continue if clean
   - updated the userSignup handler to insert the form details, return an error if any, add a flash message confirming the sign up and redirect the user to login page
 
 - 11.4 User login
   - add NonFieldErrors slice field to our validator struct for generic errors such as login failed
   - create an AddNonFieldError() method to add error messages to the NonFieldErrors slice field
   - update the valid() method to check the NonFieldErrors slice is empty
   - created and updated a login.html file
   - created a userLoginForm struct and updated the userLogin method to parse the login.html file
   - updated the Authenticate method to check if a user has valid credentials
   - updated userLoginPost method to log in a valid user and return errors for invalid users
 
 - 11.5 User logout
   - removed the session ID and value and displayed a flash message: successfully logged out
  
 - 11.6 User authorization
   - created an isAuthenticated method to return a boolean value if a user is logged in or not
   - added an IsAuthenticated boolean field to the templateData struct
   - updated the nav file to show specific info if isAuthenticated()
   - added middleware to redirect an unauthenticated user who visits /snippet/create. Set the \"Cache-Control: no-store\" so pages that require authentication are not stored in users browser cache
   - updated the routes that need to be protected

 - 11.7 CSRF (Cross-Site Request Forgery) protection
   - use of \`justinas/nosurf\` package for token-based mitigation to handle CSRF risk
   - created a noSurf handler to handle token mitigation
   - use the \`nosurf.Token()\` to get CSRF token and add to the hidden field in our forms
   - updated the forms by adding the CSRF_token field

## Chapter 12 - Using Request Context:
 - 12.1 How request context works
   - Context() can be used to store information during the lifetime of a request
   - For good practice, we create our own custom type for our context keys
  
 - 12.2 Request context for authentication/authorization
   - updated \`UserModel.Exists()\` method to return boolean if a user with specific ID exists in our \`users\` table or not
   - created context.go to define a custom context key and is authenticated context key types which provides unique key to store and retrieve authentication status from a request context
   - created an \`authenticate()\` middleware method to retrieve user's ID from session's data, validate if the user is in the database and update the request context to include the isAuthenticatedContextKey with value true
   - updated routes.go to include the authenticate() middleware in the dynamic middleware chain
   - updated isAuthenticated() helper to check the request context to determine if a user is authenticated or not, instead of checking the session data

## Chapter 13 - Optional Go features:
 - 13.1 Using embedded files
   - embedded the \`ui\` package containing our HTML, CSS, Images and JS files
   - updated routes.go to serve our CSS, JS and Images from the embedded file system, instead from disk at runtime
   - updated templates.go so our template cache uses the embeded HTML, instead of reading from disk

 - 13.2 Using generics (parametric polymorphism)
   - When to consider using generics:
     - Writing repeated boilerplate code for different data types e.g. common operations on slices, maps or channels, or helpers for carrying out validation checks or test assertions on different data types.
     - reaching for the any (empty interface{}) type. Example: creating a data structure (queue, cache or linked list) which needs to operate on different types.
   - when NOT to use generics:
     - if it makes the coder harder to understand or less clear.
     - if all the types needed to work with have a common set of methods - best to define and use a normal interface type instead.
     - just because you can

   - Writing tests in the next chapter would make us write a lot of duplicate boilerplate code.
  
   - converted PermittedInt() in validator.go to a generic function - to not only check int values, but other set of allowed values (string, int, float64 or any other comparable type).
   - updated snippetCreatePost() method handler to use the PermittedValue() in validation checks