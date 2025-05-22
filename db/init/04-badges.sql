insert into badges
(name, description)
values
(
    'Top Contributor',
    'Active in replies to topics'
),
(
    'Insightful',
    'High up/down vote ratio'
);

with vars as (
    select
        (
            select id from users
            where email = 'sample.user1@example.com'
        ) as user_id_one,
        (
            select id from users
            where email = 'another.sample.user2@example.com'
        ) as user_id_two
)

insert into user_badges
(user_id, badge_id)
values
(
    (select user_id_one from vars),
    (
        select id from badges
        where badges.name = 'Top Contributor'
    )
),
(
    (select user_id_two from vars),
    (
        select id from badges
        where badges.name = 'Insightful'
    )
);
