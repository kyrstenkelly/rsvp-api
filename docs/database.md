# Database Structure

## Event
  
| property | type      | required | description                      |
|----------|-----------|----------|----------------------------------|
| id       | INTEGER   | true     | ID of the event                  |
| name     | STRING    | true     | name of the event                |
| date     | DATE_TIME | true     | date and time of the event       |
| address_id  | INTEGER   | false    | ID of the `address` for this event |
| food_options | STRING[] | false | food options for the event | 

## Guest
| property | type     | required | description                      |
|----------|----------|----------|----------------------------------|
| id       | INTEGER  | true     | ID of the guest                  |
| name     | STRING   | true     | names of the guest               |

## Invitation
| property | type     | required | description                      |
|----------|----------|----------|----------------------------------|
| id       | INTEGER  | true     | ID of the invitation                  |
| email    | STRING   | false    | email for the invitation              |
| address_id  | INTEGER  | false    | ID of the `address` for this invitation |
| name      | STRING | true | name for the invitation (i.e. "Kelly Family") |
| plus_one | BOOLEAN  | false    | invitation includes a +1 - defaults to false |

## Invitation Guests
| property | type     | required | description                      |
|----------|----------|----------|----------------------------------|
| invitation_id | INTEGER | true | ID of the `invitation` |
| guest_id | INTEGER | true | ID of the `guest` |

## RSVP  
| property   | type    | required | description                         |
|------------|---------|----------|-------------------------------------|
| id         | INTEGER | true     | ID of the RSVP                     |
| invitation_id | INTEGER | true | ID of the `invitation` for this RSVP |
| guest_id | INTEGER | true | ID of the `guest` for this RSVP |
| attending | BOOLEAN | true | guest is coming to the event |
| food_option | STRING | false | food choice for the guest | 

## Address
| property | type    | required | description                |
|----------|---------|----------|----------------------------|
| id       | INTEGER | true     | ID of the guest            |
| line1    | STRING  | true     | first line of the address  |
| line2    | STRING  | false    | second line of the address |
| city     | STRING  | true     | city of the address        |
| state    | STRING  | true     | state of the address       |
| zip      | STRING  | true     | zip code of the address    |
