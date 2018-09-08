# Go RSVP API
Endpoint for creating RSVPs, guests and events written in Golang.

## Endpoints

### Events

* GET `/events`
* GET `/events/:event_id`
* POST `/events`
* PUT `/events/:event_id`
* DELETE `/events/:event_id`


### Guests

* GET `/guests`
* GET `/guests/:guest_id`
* POST `/guests`
* PUT `/guests/:guest_id`
* DELETE `/guests/:guest_id`

### RSVPs

* GET `/rsvps`
* GET `/rsvps/:rsvp_id`
* POST `/rsvps`
* PUT `/rsvps/:rsvp_id`
* DELETE `/rsvps/:rsvp_id`

## Models

### Event
  
| property | type      | required | description                      |
|----------|-----------|----------|----------------------------------|
| id       | INTEGER   | true     | ID of the event                  |
| name     | STRING    | true     | name of the event                |
| date     | DATE_TIME | true     | date and time of the event       |
| address  | INTEGER   | false    | ID of the `address` for this event |

### Guest
  
| property | type    | required | description                      |
|----------|---------|----------|----------------------------------|
| id       | INTEGER | true     | ID of the guest                  |
| name     | STRING  | true     | name of the guest                |
| email    | STRING  | false    | email for the guest              |
| address  | INTEGER | false    | ID of the `address` for this guest |

### RSVP
  
| property   | type    | required | description                         |
|------------|---------|----------|-------------------------------------|
| id         | INTEGER | true     | ID of the guest                     |
| event      | INTEGER | true     | ID of the `event` for this RSVP     |
| num_people | INTEGER | true     | number of people in the reservation |

### Address
  
| property | type    | required | description                |
|----------|---------|----------|----------------------------|
| id       | INTEGER | true     | ID of the guest            |
| line1    | STRING  | true     | first line of the address  |
| line2    | STRING  | false    | second line of the address |
| city     | STRING  | true     | city of the address        |
| state    | STRING  | true     | state of the address       |
| zip      | STRING  | true     | zip code of the address    |
