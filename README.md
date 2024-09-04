# Redbus Backend API Collection

Welcome to the Redbus Backend API Collection by me. This API collection powers the Redbus platform, enabling seamless bus booking and management. Below is a list of all the endpoints available in this collection along with their descriptions.

## Features
- **Bus Search**: Retrieve routes, schedules, and availability.
- **Booking Management**: Create and cancel reservations.
- **User Accounts**: Handle registration, login, forgot password and authentication.
- **Payment Processing**: Manage transactions and refunds.

### Prerequisites

- Node.js (>= 14.x)
- npm (>= 6.x)
- MongoDB

## Table of Contents

- [Register User](#register-user-endpoint)
- [Login](#login-endpoint)
- [Logout](#logout-endpoint)
- [Forgot Password](#forgot-password-endpoint)
- [Reset Password](#reset-password-endpoint)
- [Add Money to Wallet](#add-money-to-wallet)
- [Withdraw Money from Wallet](#withdraw-money-from-wallet)
- [Get Wallet Balance](#get-wallet-balance-endpoint)
- [Add New Bus](#add-new-bus-endpoint)
- [Bus Search](#bus-search-endpoint)
- [View Seats](#view-seats-endpoint)
- [Select Seat](#select-seat-endpoint)
- [Apply Offer](#apply-offer-endpoint)
- [Book Seat](#book-seat-endpoint)
- [Get Offers](#get-offers-endpoint)
- [Add Offer](#add-offer-endpoint)
- [Add Review](#add-review-endpoint)
- [Get All Reviews](#get-all-reviews-by-bus-id-endpoint)
- [Cancel Ticket](#cancel-ticket-endpoint)
- [Get User Data](#get-user-data-endpoint)
- [Get Bus by ID](#get-bus-by-id-endpoint)
- [Get Tickets by User ID](#get-tickets-by-userid-endpoint)


## Register User Endpoint

This API endpoint takes in new user details to register new user on the website
### Request

- **Method**: `POST`
- **URL**: `http://localhost:3000/register`

### Request Body
The request body should be in JSON format and include the following fields:

```json
{
    "firstName":"Check",
    "lastName":"User",
    "age":24,
    "dob":"2000-12-18",
    "email":"user@gmail.com",
    "gender":"F",
    "role":"customer",
    "encryptedPassword":"user@123",
    "profilePicPath":"/assets/pic.jpg"
}
```
## Login Endpoint

This API endpoint is used for user authentication, allowing registered users to log in to the system. It validates the user's credentials and returns a response indicating success or failure.

### Request

- **Method**: `POST`
- **URL**: `http://localhost:3000/login`

### Request Body

The request body should be in JSON format and include the following fields:

```json
{
    "email": "dhyey@gmail.com",
    "password": "dhyey@1234"
}
```
## Logout Endpoint

This API endpoints logs out the user based on token stored as cookie

### Request

- **Method**: `GET`
- **URL**: `http://localhost:3000/logout`

## Forgot Password Endpoint

This API endpoint takes email in body and checks if it exists in our database. If so, it sends a verification string to the user for reseting password

### Request

- **Method**: `POST`
- **URL**: `http://localhost:3000/forgotpassword`

### Request Body

The request body should be in JSON format and include the following fields:

```json
{
    "email":"dhyey@gmail.com"
}
```
## Reset Password Endpoint

This API endpoint resets password by verifying the key given by the user and the one temporary in the database.

### Request

- **Method**: `POST`
- **URL**: `http://localhost:3000/resetpassword`

### Request Body

The request body should be in JSON format and include the following fields:

```json
{
    "userID": "66b5e2b065ac6905744e9278",
    "key": "cKv2r4AE4rjHSkIlUJ3o1DPluNQhTLuT",
    "newPassword":"dhyey@123"
}
```
## Add Money To Wallet Endpoint

This API endpoint allows user to add money to their wallet. There is a cap to the amount a user can add at a time and also to the amount a user can have in the wallet. Only logged in users can add money

### Request

- **Method**: `POST`
- **URL**: `http://localhost:3000/addMoney/{userID}`

### Request Body

The request body should be in JSON format and include the following fields:

```json
{
    "moneyToAdd":40000
}
```
## Get Wallet Balance Endpoint

This API endpoint returns the available balance in a user's wallet

### Request

- **Method**: `GET`
- **URL**: `http://localhost:3000/getWalletBalance/{userID}`


## Withdraw Money From Wallet Endpoint

This API endpoint allows user to withdraw a certain amount from his/her wallet

### Request

- **Method**: `POST`
- **URL**: `http://localhost:3000/withdrawMoney/{userID}`

### Request Body

The request body should be in JSON format and include the following fields:

```json
{
    "moneyToWithdraw":300
}
```
## Add New Bus Endpoint

This request takes in the start location, the end location and the date of your travel. It also takes in filters as a map where keys are strings and values are bool.

### Request

- **Method**: `POST`
- **URL**: `http://localhost:3000/addBus`

### Request Body

The request body should be sent as form-data. Below is a sample payload:

| Key                             | Value                                               | Type   |
|---------------------------------|-----------------------------------------------------|--------|
| `operatorName`                   | New Motors                                          | text   |
| `modelDetails`                   | Tata                                                | text   |
| `totalSeats`                     | 45                                                  | text   |
| `imgPath`                        | [file] (upload image file)                         | file   |
| `amenities`                      | chargingPoint, wifi, toilet, blankets, waterBottle | text   |
| `avgRating`                      | 0                                                   | text   |
| `numberOfReviews`                | 0                                                   | text   |
| `liveTracking`                   | true                                                | text   |
| `isAcAvailable`                  | true                                                | text   |
| `busType`                        | 30SL                                                | text   |
| `frequency`                      | Weekends                                            | text   |
| `seatAvailability`              | 45                                                  | text   |
| `sleeperCost`                    | 1.9                                                 | text   |
| `seaterCost`                     | 1.2                                                 | text   |
| `stops[0][location]`             | Bengaluru                                           | text   |
| `stops[0][arrivalTime]`          |                                                     | text   |
| `stops[0][departureTime]`        | 10:00                                               | text   |
| `stops[1][location]`             | Mysuru                                              | text   |
| `stops[1][arrivalTime]`          | 12:30                                               | text   |
| `stops[1][departureTime]`        | 12:40                                               | text   |
| `stops[2][location]`             | Ooty                                                | text   |
| `stops[2][arrivalTime]`          | 19:40                                               | text   |
| `stops[2][departureTime]`        |                                                     | text   |
| `stops[3][location]`             | Chennai                                             | text   |
| `stops[3][arrivalTime]`          | 18:20                                               | text   |
| `stops[3][departureTime]`        |                                                     | text   |

## Bus Search Endpoint

This API endpoint searches bus on the basis of travel start, end and dates. It returns an array of booking objects which provides with the info of seat availability, travel dates and times. Pagination is also applied

### Request

- **Method**: `GET`
- **URL**: `http://localhost:3000/buses/search/{page}`

### Request Body

The request body should be in JSON format and include the following fields:

```json
{
    "fromLocation":"Ooty",
    "toLocation":"Vellore",
    "travelDate":"2024-08-24"
}
```
## View Seats Endpoint
This API endpoint uses gets the seating arrangement, availability and pricing based on the booking

### Request

- **Method**: `GET`
- **URL**: `http://localhost:3000/viewseats/{busID}`

### Request Body

The request body should be in JSON format and include the following fields:

```json
{
    "fromLocation":"Ooty",
    "toLocation":"Vellore",
    "travelDate":"2024-08-24"
}
```
## Select Seat Endpoint
This API endpoint takes the seatID and bus ID to give you the price of the seat selected. This also takes into consideration if the seat is selected or not by the user. A frontend engineer can directly use this data to modify the total price shown to the user 

### Request

- **Method**: `GET`
- **URL**: `http://localhost:3000/selectseat/{seatID}/{busID}`

### Request Body

The request body should be in JSON format and include the following fields:

```json
{
    "fromLocation":"Ooty",
    "toLocation":"Vellore",
    "travelDate":"2024-08-24"
}
```
## Apply Offer Endpoint

This API endpoint takes in the offer code and the total cart value including taxes and returns the final price after offer. It also checks if the offer is valid by date and by the total cart value

### Request

- **Method**: `POST`
- **URL**: `http://localhost:3000/applyOffer`

### Request Body

The request body should be in JSON format and include the following fields:

```json
{
    "cartValue":867,
    "offerCode":"FESTIVE30"
}
```
## Book Seat Endpoint

This API endpoint uses the booking and payment details to book tickets for a user. Gender specificity is taken care of while booking seats. Also wallet balances are also checked

### Request

- **Method**: `POST`
- **URL**: `http://localhost:3000/bookseat`

### Request Body

The request body should be in JSON format and include the following fields:

```json
{"booking":{
        "busID": "66c2359890dd89a96015583e",
        "travelStartDate": "2024-08-24",
        "travelEndDate": "2024-08-24",
        "travelStartTime": "11:10",
        "travelEndTime": "16:10",
        "travelStartLocation": "Ooty",
        "travelEndLocation": "Vellore"
    },
    "seatIDs":["66c2359890dd89a96015580d"],
    "passengerNames":["Roharma"],
    "passengerGenders":["M"],
    "passengerAges":[21],
    "baseFare": 867,
    "discountedAmount": 273,
    "gst": 43,
    "totalPayableAmount": 637
    }
```
## Get Offers Endpoint

This API endpoint returns an array of all Offer objects

### Request

- **Method**: `GET`
- **URL**: `http://localhost:3000/getOffers`

## Add Offer Endpoint

This API endpoint allows the admin to add offer. It takes offer description, code, validity and other necessary details in body

### Request

- **Method**: `POST`
- **URL**: `http://localhost:3000/addOffer`

### Request Body

The request body should be in JSON format and include the following fields:

```json
{
    "oCode":"FESTIVE30",
    "description":"Flat 30% off upto Rs.400 on orders above Rs.3000",
    "validity":"2024-08-30",
    "minOrderVal":3000,
    "maxDiscount":400,
    "discount":30
}
```
## Add Review Endpoint

This API endpoint takes in the user ID, and the bus ID to review. It will also check for user login. It will also take into consideration if the user can rate the bus based on whether he/she has completed the travel

### Request

- **Method**: `POST`
- **URL**: `http://localhost:3000/addReview/{userID}/{busID}`

### Request Body

The request body should be in JSON format and include the following fields:

```json
{
    "rating":3,
    "reviewText":"This is a good bus"
}
```
## Get All Reviews By Bus ID Endpoint

This API endpoint fetches all reviews based on the bus ID

### Request

- **Method**: `GET`
- **URL**: `http://localhost:3000/getReviews/{busID}`

## Cancel Ticket Endpoint

This API endpoint cancels ticket by ticket ID and refunds amount based on the time before which the user has applied for cancellation and refund.

### Request

- **Method**: `DELETE`
- **URL**: `http://localhost:3000/cancelTicket/{ticketID}`

### Request Body

The request body should be in JSON format and include the following fields:

```json
{
    "email":"jane@gmail.com",
    "password":"jane@123"
}
```
## Get User Data Endpoint

This API endpoint returns all details of a user

### Request

- **Method**: `GET`
- **URL**: `http://localhost:3000/users/{userID}`

## Get Tickets By UserID Endpoint

This API endpoint returns an array of all tickets booked by a user
### Request

- **Method**: `GET`
- **URL**: `http://localhost:3000/getTickets/{userID}`


## Get Bus By ID Endpoint

This API endpoint returns complete bus details based on the bus ID

### Request

- **Method**: `GET`
- **URL**: `http://localhost:3000/buses/{busID}`
