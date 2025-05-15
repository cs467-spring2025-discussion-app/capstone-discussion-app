do $$
begin
  for i in 1..12 loop
    insert into discussion_topics
    (title, prompt, is_active)
    values
    (
      'sample topic ' || i,
      'description ' || i,
      true
    );
  end loop;
end $$;

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

insert into replies
(discussion_id, user_id, content)
values
(1, (select user_id_one from vars), 'reply 1-1'),
(1, (select user_id_two from vars), 'reply 1-2'),
(2, (select user_id_one from vars), 'reply 2-1'),
(3, (select user_id_two from vars), 'reply 3-1'),
(3, (select user_id_one from vars), 'reply 3-2'),
(3, (select user_id_one from vars), 'reply 3-3'),
(4, (select user_id_two from vars), 'reply 4-1');