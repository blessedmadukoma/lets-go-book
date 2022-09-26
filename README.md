1. Added bcrypt package for hashing passwords
2. updated insert method to hash password, check if email is duplicate and either return an error if there is any or continue if clean
3. updated the userSignup handler to insert the form details, return an error if any, add a flash message confirming the sign up and redirect the user to login page