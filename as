SELECT p.recruiter_public_id, COUNT(*) as interview_count
FROM interviews i
INNER JOIN user_interviews ui ON ui.interview_id = i.id
INNER JOIN positions p ON p.id = ui.position_id
INNER JOIN recruiters r ON p.recruiter_public_id = r.public_id
WHERE p.recruiter_public_id = '48c1d92d-6991-468f-8b30-fff1179d9b2a'
GROUP BY p.recruiter_public_id;