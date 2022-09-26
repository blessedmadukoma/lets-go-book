#### Let's Go Code Implementation - Alex Edwards
## Chapter 11:
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