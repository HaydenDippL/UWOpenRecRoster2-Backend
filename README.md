# UWOpenRecRoster2-Backend

Backend written in GO for [UWOpenRecRoster2](https://github.com/HaydenDippL/UWOpenRecRoster2). Check our DOCS.md for information about the code in the project.

## Endpoint

`GET /schedule` with the parameters:

- `date` with a date in the ISO Date format `"2024-12-25` for December 25, 2024
- `gyms` represents which gyms we want. This paramter is a comma separated list which may contain `"Bakke"` and or `"Nick"`
- `facilities` represents what type of facilities we are are querying. This is a comma separated list, similar to `gym` with the following possible options: `"Courts"`, `"Pool"`, `"RockWall"`, `"Esports"`.

A full query may look like 

`GET /schedule?date=12-25-2025&gym=Bakke,Nick&facilities=Courts,Pool`

It returns data in the following form, all dates are ISO DateTimes (RFC 3339).

```
{
    Bakke: {
        Courts: [
            {
                location: "Court 1",
                eventName: "Open Rec Basketball",
                start: "2025-03-11T06:00:00Z",
                end: "2025-03-11T10:00:00Z"
            }
            ...
        ],
        Pool: [
            {
                location: "Lane 1",
                eventName: "Open Rec Swim",
                start: "2025-03-11T06:00:00Z",
                end: "2025-03-11T10:00:00Z"
            }
        ]
    },
    Nick: {
        Courts: [
            ...
        ],
        Pool: [
            ...
        ]
    }
}
```

## Requirements

Need Go and required build tools

```
sudo apt-get install -y \
    build-essential \
    gcc-multilib \
    linux-libc-dev \
    libgcc-s1 \
    libc6-dev
```

Install postgres

```
sudo apt install postgresql postgresql-contrib
sudo systemctl start postgresql
sudo systemctl enable postgresql
```

```
go get github.com/gin-gonic/gin
go get gorm.io/gorm
go get gorm.io/driver/postgres
```

## PSQL

Connect to Postgres DB

```
psql -h <DB_HOST> -U <DB_USER> -d <DB_NAME> -p <DB_PORT>
```