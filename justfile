# justfile docs     https://just.systems/man/en/
# This file is used for local development convenience
# Dependencies:
# - just            https://github.com/casey/just
# - psql            https://archlinux.org/packages/?name=postgresql
# - docker-compose  https://docs.docker.com/compose/install/
# - curl            https://curl.se/
# - jq              https://github.com/jqlang/jq

set quiet := true

FAST := "localhost:8080"
AUTH := "localhost:3001"
DEFAULT_USER := "bob@just.com"
DEFAULT_PASS := "s+r0nGpw0rd"

default:
    @just --list

# #############################################################################
# API endpoint recipes
# #############################################################################
# ##################
# GoAuth
# ##################

# goauth/register
register email=(DEFAULT_USER) pw=(DEFAULT_PASS):
    just request {{ AUTH }} register \
    POST '{\"email\":\"{{ email }}\", \"password\":\"{{ pw }}\"}' | jq

# goauth/login
login email=(DEFAULT_USER) pw=(DEFAULT_PASS):
    curl -s --request POST \
        --data '{"email":"{{ email }}", "password":"{{ pw }}"}' \
        --header 'Content-Type: application/json' \
        -c /tmp/go_capstone_auth_{{ email}}_cookies.txt \
        {{ AUTH }}/login | jq

whoami email=(DEFAULT_USER):
    curl -s --request GET \
        -b /tmp/go_capstone_auth_{{ email}}_cookies.txt \
        {{ AUTH }}/whoami | jq


# goauth/logout
logout email=(DEFAULT_USER):
    curl -s --request POST \
        --header 'Content-Type: application/json' \
        -b /tmp/go_capstone_auth_{{ email}}_cookies.txt \
        {{ AUTH }}/logout | jq

# goauth/logouteverywhere
logouteverywhere email=(DEFAULT_USER):
    curl -s --request POST \
        --header 'Content-Type: application/json' \
        -b /tmp/go_capstone_auth_{{ email}}_cookies.txt \
        {{ AUTH }}/logouteverywhere | jq

# goauth/updateuser
updateuser email=(DEFAULT_USER) pw=(DEFAULT_PASS):
    curl -s --request POST \
        --data '{"email":"{{ email }}", "password":"{{ pw }}"}' \
        --header 'Content-Type: application/json' \
        -b /tmp/go_capstone_auth_{{ email}}_cookies.txt \
        {{ AUTH }}/updateuser | jq

# goauth/deleteaccount
deleteaccount email=(DEFAULT_USER):
    curl -s --request DELETE \
        -b /tmp/go_capstone_auth_{{ email}}_cookies.txt \
        {{ AUTH }}/deleteaccount | jq



# ##################
# FastApi
# ######
# ######
# Topics
# ######

# fastapi/topics/submit_topic
submit_topic title prompt:
    just request {{ FAST }} topics/submit_topic \
    POST '{\"title\":\"{{ title }}\", \"prompt\":\"{{ prompt }}\"}' | jq

# fastapi/topics/get_daily_topics
get_daily_topics:
    just request {{ FAST }} topics/get_daily_topics GET | jq

# fastapi/topics/get_topic/{topic_id}
get_topic id:
    just request {{ FAST }} topics/get_topic/{{ id }} GET | jq

# fastapi/topics/get_topic/{topic_id}
get_all_topics:
    just request {{ FAST }} topics/get_all_topics GET | jq

# fastapi/topics/activate_topic/{topic_id}
activate_topic id:
    just request {{ FAST }} topics/activate_topic/{{ id }} PUT | jq

# fastapi/topics/deactivate_topic/{topic_id}
deactivate_topic id:
    just request {{ FAST }} topics/deactivate_topic/{{ id }} PUT | jq

# fastapi/topics/show_topic/{topic_id}
show_topic id:
    just request {{ FAST }} topics/show_topic/{{ id }} PUT | jq

# fastapi/topics/update_topic/{topic_id}
update_topic id title prompt:
    just request {{ FAST }} topics/update_topic/{{ id }} \
    PUT '{\"title\":\"{{ title }}\", \"prompt\":\"{{ prompt }}\"}' | jq

# fastapi/topics/delete_topic/{topic_id}
delete_topic id:
    just request {{ FAST }} topics/delete_topic/{{ id }} DELETE | jq

rotate:
    just request {{ FAST }} topics/rotate PUT | jq

# ######
# Replies
# ######

# fastapi/replies/submit_reply/{discussion_id}
submit_reply email discussion_id content:
    #!/usr/bin/env sh
    query="select id from users where email='{{ email }}'"
    user_id=$(docker exec app_db psql \
        -U postgresdiscussionapp -d app_db_discussion \
        -t -A \
        -c "$query")

    just request {{ FAST }} replies/submit_reply/{{ discussion_id }} \
    POST '{\"user_id\":\"'$user_id'\", \"content\":\"{{ content }}\"}' | jq

# fastapi/replies/get_all_replies/{discussion_id}
get_all_replies discussion_id:
    just request {{ FAST }} replies/get_all_replies/{{ discussion_id }} GET | jq

# fastapi/replies/get_reply/{reply_id}
get_reply reply_id:
    just request {{ FAST }} replies/get_reply/{{ reply_id }} GET | jq

# fastapi/replies/counts
reply_counts:
    just request {{ FAST }} replies/counts GET | jq

# fastapi/replies/update_reply/{reply_id}
update_reply email reply_id content:
    #!/usr/bin/env sh
    query="select id from users where email='{{ email }}'"
    user_id=$(docker exec app_db psql \
        -U postgresdiscussionapp -d app_db_discussion \
        -t -A \
        -c "$query")
    just request {{ FAST }} replies/update_reply/{{ reply_id }} \
    PUT '{\"user_id\":\"'$user_id'\", \"content\":\"{{ content }}\"}' | jq

# fastapi/topics/delete_reply/{topic_id}
delete_reply id:
    just request {{ FAST }} replies/delete_reply/{{ id }} POST | jq

# ######
# Votes
# ######

# fastapi/votes/submit_vote/{reply_id}
vote email=(DEFAULT_USER) reply_id="" type="0":
    #!/usr/bin/env sh
    query="select id from users where email='{{ email }}'"
    user_id=$(docker exec app_db psql \
        -U postgresdiscussionapp -d app_db_discussion \
        -t -A \
        -c "$query")

    just request {{ FAST }} votes/submit_vote/{{ reply_id }} \
    POST '{\"user_id\":\"'$user_id'\", \"vote_type\": {{ type }} }' | jq

# #############################################################################
# Direct database query recipes
# #############################################################################
# ######
# Users
# ######

# Get all registered users
users:
    just exec "select * from users"

# Get a user id by email
id email:
    just exec "select id from users where email='{{ email }}'" "-t -A"

# ######
# Topics
# ######

# Get all topics and their reply counts
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

# #############################################################################
# Docker recipes
# #############################################################################

# Stop the app_db service and delete the app_postgres_data volume
empty:
    docker compose down db
    docker volume rm capstone-discussion-app_postgres_data

# Stop and rebuild the docker services
rebuild:
    docker compose down
    docker compose build --no-cache
    docker compose up

# Execute a psql command in the docker database
exec command flags="":
    docker exec -it app_db psql \
        -U postgresdiscussionapp -d app_db_discussion \
        {{ flags }} \
        --command "{{ command }}"

# #############################################################################
# Misc recipes
# #############################################################################

# Make an http request
request server endpoint method data="":
    curl -s --request {{ method }} \
        --data "{{ data }}" \
        --header 'Content-Type: application/json' {{ server }}/{{ endpoint }}

# Ping either `auth` or `fast`
ping service:
    #!/usr/bin/env sh
    if [ "{{ service }}" = "auth" ]; then
        curl {{ AUTH }}/ping
    elif [ "{{ service }}" = "fast" ]; then
        curl {{ FAST }}
    fi
    printf "\n"
