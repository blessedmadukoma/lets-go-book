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
   
   **Verifying the user details**
   - updated the Authenticate method to check if a user has valid credentials
   - updated userLoginPost method to log in a valid user and return errors for invalid users