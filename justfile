# justfile docs     https://just.systems/man/en/
# This file is used for local development convenience
# Dependencies:
# - just            https://github.com/casey/just
# - psql            https://archlinux.org/packages/?name=postgresql
# - docker-compose  https://docs.docker.com/compose/install/
# - curl            https://curl.se/
# - jq              https://github.com/jqlang/jq

set quiet := true

db_user := "postgres_capstone_app"
db_name := "capstone_app"
db_volume := "db"
FAST := "localhost:8080"
AUTH := "localhost:3001"
def_user := "bob@just.com"
def_pw := "s+r0nGpw0rd"

default:
    @just --list

# #############################################################################
# API endpoint recipes
# #############################################################################
# ##################
# GoAuth
# ##################

# `goauth/register`
[group('authapi')]
register email=(def_user) pw=(def_pw):
    just request {{ AUTH }} register \
    POST '{\"email\":\"{{ email }}\", \"password\":\"{{ pw }}\"}' | jq

# `goauth/login`
[group('authapi')]
login email=(def_user) pw=(def_pw):
    curl -s --request POST \
        --data '{"email":"{{ email }}", "password":"{{ pw }}"}' \
        --header 'Content-Type: application/json' \
        -c /tmp/go_capstone_auth_{{ email }}_cookies.txt \
        {{ AUTH }}/login | jq

# `goauth/whoami`
[group('authapi')]
whoami email=(def_user):
    curl -s --request GET \
        -b /tmp/go_capstone_auth_{{ email }}_cookies.txt \
        {{ AUTH }}/whoami | jq

# `goauth/logout`
[group('authapi')]
logout email=(def_user):
    curl -s --request POST \
        --header 'Content-Type: application/json' \
        -b /tmp/go_capstone_auth_{{ email }}_cookies.txt \
        {{ AUTH }}/logout | jq

# `goauth/logouteverywhere`
[group('authapi')]
logouteverywhere email=(def_user):
    curl -s --request POST \
        --header 'Content-Type: application/json' \
        -b /tmp/go_capstone_auth_{{ email }}_cookies.txt \
        {{ AUTH }}/logouteverywhere | jq

# `goauth/updateuser`
[group('authapi')]
updateuser email=(def_user) pw=(def_pw):
    curl -s --request POST \
        --data '{"email":"{{ email }}", "password":"{{ pw }}"}' \
        --header 'Content-Type: application/json' \
        -b /tmp/go_capstone_auth_{{ email }}_cookies.txt \
        {{ AUTH }}/updateuser | jq

# `goauth/deleteaccount`
[group('authapi')]
deleteaccount email=(def_user):
    curl -s --request DELETE \
        -b /tmp/go_capstone_auth_{{ email }}_cookies.txt \
        {{ AUTH }}/deleteaccount | jq

# ##################
# FastApi
# ######
# ######
# Topics
# ######

# `fastapi/topics/submit_topic`
[group('fastapi/topics')]
submit_topic title prompt:
    just request {{ FAST }} topics/submit_topic \
    POST '{\"title\":\"{{ title }}\", \"prompt\":\"{{ prompt }}\"}' | jq

# `fastapi/topics/get_daily_topics`
[group('fastapi/topics')]
get_daily_topics:
    just request {{ FAST }} topics/get_daily_topics GET | jq

# `fastapi/topics/get_topic/{topic_id}`
[group('fastapi/topics')]
get_topic id:
    just request {{ FAST }} topics/get_topic/{{ id }} GET | jq

# `fastapi/topics/get_topic/{topic_id}`
[group('fastapi/topics')]
get_all_topics:
    just request {{ FAST }} topics/get_all_topics GET | jq

# `fastapi/topics/activate_topic/{topic_id}`
[group('fastapi/topics')]
activate_topic id:
    just request {{ FAST }} topics/activate_topic/{{ id }} PUT | jq

# `fastapi/topics/deactivate_topic/{topic_id}`
[group('fastapi/topics')]
deactivate_topic id:
    just request {{ FAST }} topics/deactivate_topic/{{ id }} PUT | jq

# `fastapi/topics/show_topic/{topic_id}`
[group('fastapi/topics')]
show_topic id:
    just request {{ FAST }} topics/show_topic/{{ id }} PUT | jq

# `fastapi/topics/update_topic/{topic_id}`
[group('fastapi/topics')]
update_topic id title prompt:
    just request {{ FAST }} topics/update_topic/{{ id }} \
    PUT '{\"title\":\"{{ title }}\", \"prompt\":\"{{ prompt }}\"}' | jq

# `fastapi/topics/delete_topic/{topic_id}`
[group('fastapi/topics')]
delete_topic id:
    just request {{ FAST }} topics/delete_topic/{{ id }} DELETE | jq

# `fastapi/topics/rotate`
[group('fastapi/topics')]
rotate:
    just request {{ FAST }} topics/rotate PUT | jq

# ######
# Replies
# ######

# `fastapi/replies/submit_reply/{discussion_id}`
[group('fastapi/replies')]
submit_reply email discussion_id content:
    #!/usr/bin/env sh
    query="select id from users where email='{{ email }}'"
    user_id=$(docker exec app_db psql \
        -U {{ db_user }} -d {{ db_name }} \
        -t -A \
        -c "$query")

    just request {{ FAST }} replies/submit_reply/{{ discussion_id }} \
    POST '{\"user_id\":\"'$user_id'\", \"content\":\"{{ content }}\"}' | jq

# `fastapi/replies/get_all_replies/{discussion_id}`
[group('fastapi/replies')]
get_all_replies discussion_id:
    just request {{ FAST }} replies/get_all_replies/{{ discussion_id }} GET | jq

# `fastapi/replies/get_reply/{reply_id}`
[group('fastapi/replies')]
get_reply reply_id:
    just request {{ FAST }} replies/get_reply/{{ reply_id }} GET | jq

# `fastapi/replies/counts`
[group('fastapi/replies')]
reply_counts:
    just request {{ FAST }} replies/counts GET | jq

# `fastapi/replies/update_reply/{reply_id}`
[group('fastapi/replies')]
update_reply email reply_id content:
    #!/usr/bin/env sh
    query="select id from users where email='{{ email }}'"
    user_id=$(docker exec app_db psql \
        -U {{ db_user }} -d {{ db_name }} \
        -t -A \
        -c "$query")
    just request {{ FAST }} replies/update_reply/{{ reply_id }} \
    PUT '{\"user_id\":\"'$user_id'\", \"content\":\"{{ content }}\"}' | jq

# `fastapi/topics/delete_reply/{topic_id}`
[group('fastapi/replies')]
delete_reply id:
    just request {{ FAST }} replies/delete_reply/{{ id }} POST | jq

# ######
# Votes
# ######

# `fastapi/votes/submit_vote/{reply_id}`
[group('fastapi/votes')]
vote email=(def_user) reply_id="" type="0":
    #!/usr/bin/env sh
    query="select id from users where email='{{ email }}'"
    user_id=$(docker exec app_db psql \
        -U {{ db_user }} -d {{ db_name }} \
        -t -A \
        -c "$query")

    just request {{ FAST }} votes/submit_vote/{{ reply_id }} \
    POST '{\"user_id\":\"'$user_id'\", \"vote_type\": {{ type }} }' | jq

# Get all rows from the `votes` table
votes:
    just exec "\
        select \
            users.email, \
            votes.reply_id, \
            votes.vote_type, \
            votes.created_at, \
            votes.updated_at \
        from votes \
        inner join users on votes.user_id=id \
    "

# #############################################################################
# Direct database query recipes
# #############################################################################

# Source a sql file from the db/init dir into the app database
[group('sql')]
source file:
    docker exec -it app_db psql -U {{ db_user }} -d {{ db_name }} -f docker-entrypoint-initdb.d/{{ file }}

# ######
# Users
# ######

# Execute a psql command in the docker database
[group('sql')]
exec command flags="":
    docker exec -it app_db psql \
        -U {{ db_user }} -d {{ db_name }} \
        {{ flags }} \
        --command "{{ command }}"

# Get all registered users
[group('sql')]
users:
    just exec "select * from users"

# Get a user id by email
[group('sql')]
id email:
    just exec "select id from users where email='{{ email }}'" "-t -A"

# Lock a user
[group('sql')]
lock email time='2099-01-01 00:00:00':
    just exec "\
        update users \
        set \
            account_locked=true, \
            account_locked_until='{{ time }}' \
        where email='{{ email }}' \
    "

# Unlock a user
[group('sql')]
unlock email:
    just exec "\
        update users \
            set \
                account_locked=false, \
                account_locked_until=null, \
                failed_login_attempts=0 \
        where email='{{ email }}' \
    "

# ######
# Topics
# ######

# Get all topics and their reply counts
[group('sql')]
topics:
    just exec "\
    select \
        discussion_topics.id, \
        discussion_topics.title, \
        discussion_topics.prompt, \
        discussion_topics.is_active, \
        ( \
        select \
            count(distinct replies.id) num_replies \
        from replies \
        where replies.discussion_id = discussion_topics.id \
        ) \
    from discussion_topics"

# ######
# Badges
# ######

# Get all user badges
[group('sql')]
ubadges:
    just exec "\
    select \
        users.email, \
        badges.name, \
        user_badges.awarded_at \
    from user_badges \
    inner join users on user_badges.user_id=users.id \
    inner join badges on user_badges.badge_id=badges.id \
    "

# ######
# Reputation
# ######

# Get all user reputations listed by email
[group('sql')]
rep:
    just exec "\
    select \
        users.email, \
        user_reputation.upvotes_received, \
        user_reputation.downvotes_received \
    from user_reputation \
    inner join users on user_reputation.user_id=users.id \
    "


# #############################################################################
# Docker recipes
# #############################################################################

# Stop the app_db service and delete the database volume
[group('docker')]
empty:
    docker compose down -v {{ db_volume }}

# Rebuild services from scratch and start
[group('docker')]
fresh:
    docker compose down
    docker network prune -f
    docker compose build --no-cache
    docker compose up

# Start all services, rebuild if necessary
[group('docker')]
start *flags:
    docker compose up -d --build {{ flags }}

# Stop all services
[group('docker')]
stop:
    docker compose stop

# Show logs (optional flags e.g. `-f --timestamps`)
[group('docker')]
logs *flags:
    docker compose logs {{ flags }}

# Show logs with `follow` and `timestampts option`
[group('docker')]
flogs service='':
    just logs {{ service }} -f --timestamps

# #############################################################################
# Misc recipes
# #############################################################################

# Make an http request
[group('misc')]
request server endpoint method data="":
    curl -s --request {{ method }} \
        --data "{{ data }}" \
        --header 'Content-Type: application/json' {{ server }}/{{ endpoint }}

# Ping either `auth` or `fast`
[group('misc')]
ping service:
    #!/usr/bin/env sh
    if [ "{{ service }}" = "auth" ]; then
        curl {{ AUTH }}/ping
    elif [ "{{ service }}" = "fast" ]; then
        curl {{ FAST }}
    fi
    printf "\n"
