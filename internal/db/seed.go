package db

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"math/rand"
	"strconv"

	"github.com/gad/social/internal/store"
)

var potentialNames = []string{
	"Alice", "Bob", "Charlie", "David", "Eve", "Frank", "Grace", "Hank", "Ivy",
	"Jack", "Karen", "Leo", "Mona", "Nina", "Oscar", "Paul", "Quincy", "Rachel",
	"Steve", "Tina", "Uma", "Victor", "Wendy", "Xander", "Yara", "Zack", "Abby",
	"Ben", "Cara", "Derek", "Ella", "Felix", "Gina", "Hugo", "Isla", "Jake",
	"Kira", "Liam", "Mia", "Noah", "Olive", "Pete", "Quinn", "Riley", "Sara",
	"Tom", "Ursula", "Vera", "Will", "Xena", "Yvonne", "Zane", "Ava", "Bella",
	"Chloe", "Daisy", "Emma", "Fiona", "George", "Hannah", "Ian", "Jade",
	"Kyle", "Luna", "Molly", "Nora", "Owen", "Penny", "Quin", "Ruby",
	"Sophie", "Theo", "Uma", "Violet", "Wyatt", "Xyla", "Yara", "Zelda",
	"Aaron", "Beth", "Cody", "Diana", "Eli", "Faith", "Gabe", "Holly",
	"Ira", "Jenna", "Kai", "Lila", "Miles", "Nina", "Otis", "Piper",
	"Quinn", "Rosa", "Seth", "Tara", "Uri", "Vera", "Wade", "Xander",
	"Yara", "Zoe",
}

var potentialContents = []string{
	"Chase your dreams, but don’t forget to enjoy the journey along the way.",
	"Every small step you take brings you closer to your big goals.",
	"Life is short; make every moment count and every smile matter.",
	"Be the reason someone believes in the goodness of people today.",
	"Your vibe attracts your tribe; stay positive and surround yourself with light.",
	"Success isn’t about perfection; it’s about progress and persistence.",
	"Kindness costs nothing but means everything to someone who needs it.",
	"Dream big, work hard, stay focused, and surround yourself with good people.",
	"Sometimes the smallest step in the right direction ends up being the biggest.",
	"Don’t wait for the perfect moment; take the moment and make it perfect.",
	"Your attitude determines your direction; choose positivity and watch your life change.",
	"Be so good they can’t ignore you; let your work speak for itself.",
	"Every day is a new opportunity to grow, learn, and become better.",
	"Don’t compare your chapter 1 to someone else’s chapter 20; trust your journey.",
	"Stay true to yourself; the right people will love the real you.",
	"Hard work beats talent when talent doesn’t work hard; keep pushing forward.",
	"Life is better when you’re laughing; find joy in the little things.",
	"Believe in yourself, even when no one else does; you’ve got this.",
	"Success is the sum of small efforts repeated day in and day out.",
	"Don’t let fear of failure stop you from chasing what sets your soul on fire.",
	"Be the kind of person who makes others feel seen, heard, and valued.",
	"Your time is limited; don’t waste it living someone else’s life.",
	"Every storm runs out of rain; keep going, brighter days are ahead.",
	"Surround yourself with those who lift you higher and inspire you to grow.",
	"Don’t just exist; live boldly, love deeply, and leave a mark.",
	"Small progress is still progress; celebrate every step forward.",
	"Your potential is endless; don’t let self-doubt hold you back.",
	"Be the energy you want to attract; positivity is contagious.",
	"Life is a canvas; make sure you paint it with vibrant colors.",
	"Don’t wait for opportunities; create them with hard work and determination.",
	"Stay humble, work hard, and be kind; success will follow.",
	"Your dreams are valid; keep working toward them no matter what.",
	"Every setback is a setup for a comeback; keep pushing forward.",
	"Be the light in someone’s darkness; kindness can change everything.",
	"Don’t let yesterday take up too much of today; focus on now.",
	"Success is a journey, not a destination; enjoy every step.",
	"Your mindset determines your success; think positively and take action.",
	"Life is too short to waste on negativity; focus on what matters.",
	"Be fearless in the pursuit of what sets your soul on fire.",
	"Every day is a chance to rewrite your story; make it count.",
	"Don’t let fear stop you from becoming the best version of yourself.",
	"Your vibe attracts your tribe; stay true to who you are.",
	"Dream big, work hard, stay focused, and never give up.",
	"Small steps every day lead to big results over time; keep going.",
	"Be the reason someone smiles today; kindness is always in style.",
	"Your potential is limitless; don’t let doubt hold you back.",
	"Life is a gift; unwrap it with gratitude and joy.",
	"Stay true to yourself; the right people will love the real you.",
	"Every day is a new beginning; make it count and shine bright.",
	"Don’t wait for the perfect moment; take the moment and make it perfect.",
	"Your attitude determines your direction; choose positivity and watch your life change.",
	"Be so good they can’t ignore you; let your work speak for itself.",
	"Every day is a new opportunity to grow, learn, and become better.",
	"Don’t compare your chapter 1 to someone else’s chapter 20; trust your journey.",
	"Stay true to yourself; the right people will love the real you.",
	"Hard work beats talent when talent doesn’t work hard; keep pushing forward.",
	"Life is better when you’re laughing; find joy in the little things.",
	"Believe in yourself, even when no one else does; you’ve got this.",
	"Success is the sum of small efforts repeated day in and day out.",
	"Don’t let fear of failure stop you from chasing what sets your soul on fire.",
	"Be the kind of person who makes others feel seen, heard, and valued.",
	"Your time is limited; don’t waste it living someone else’s life.",
	"Every storm runs out of rain; keep going, brighter days are ahead.",
	"Surround yourself with those who lift you higher and inspire you to grow.",
	"Don’t just exist; live boldly, love deeply, and leave a mark.",
	"Small progress is still progress; celebrate every step forward.",
	"Your potential is endless; don’t let self-doubt hold you back.",
	"Be the energy you want to attract; positivity is contagious.",
	"Life is a canvas; make sure you paint it with vibrant colors.",
	"Don’t wait for opportunities; create them with hard work and determination.",
	"Stay humble, work hard, and be kind; success will follow.",
	"Your dreams are valid; keep working toward them no matter what.",
	"Every setback is a setup for a comeback; keep pushing forward.",
	"Be the light in someone’s darkness; kindness can change everything.",
	"Don’t let yesterday take up too much of today; focus on now.",
	"Success is a journey, not a destination; enjoy every step.",
	"Your mindset determines your success; think positively and take action.",
	"Life is too short to waste on negativity; focus on what matters.",
	"Be fearless in the pursuit of what sets your soul on fire.",
	"Every day is a chance to rewrite your story; make it count.",
	"Don’t let fear stop you from becoming the best version of yourself.",
	"Your vibe attracts your tribe; stay true to who you are.",
	"Dream big, work hard, stay focused, and never give up.",
	"Small steps every day lead to big results over time; keep going.",
	"Be the reason someone smiles today; kindness is always in style.",
	"Your potential is limitless; don’t let doubt hold you back.",
	"Life is a gift; unwrap it with gratitude and joy.",
	"Stay true to yourself; the right people will love the real you.",
	"Every day is a new beginning; make it count and shine bright.",
	"Don’t wait for the perfect moment; take the moment and make it perfect.",
	"Your attitude determines your direction; choose positivity and watch your life change.",
	"Be so good they can’t ignore you; let your work speak for itself.",
	"Every day is a new opportunity to grow, learn, and become better.",
	"Don’t compare your chapter 1 to someone else’s chapter 20; trust your journey.",
	"Stay true to yourself; the right people will love the real you.",
	"Hard work beats talent when talent doesn’t work hard; keep pushing forward.",
	"Life is better when you’re laughing; find joy in the little things.",
	"Believe in yourself, even when no one else does; you’ve got this.",
	"Success is the sum of small efforts repeated day in and day out.",
	"Don’t let fear of failure stop you from chasing what sets your soul on fire.",
	"Be the kind of person who makes others feel seen, heard, and valued.",
	"Your time is limited; don’t waste it living someone else’s life.",
	"Every storm runs out of rain; keep going, brighter days are ahead.",
	"Surround yourself with those who lift you higher and inspire you to grow.",
	"Don’t just exist; live boldly, love deeply, and leave a mark.",
	"Small progress is still progress; celebrate every step forward.",
	"Your potential is endless; don’t let self-doubt hold you back.",
	"Be the energy you want to attract; positivity is contagious.",
	"Life is a canvas; make sure you paint it with vibrant colors.",
	"Don’t wait for opportunities; create them with hard work and determination.",
	"Stay humble, work hard, and be kind; success will follow.",
	"Your dreams are valid; keep working toward them no matter what.",
	"Every setback is a setup for a comeback; keep pushing forward.",
	"Be the light in someone’s darkness; kindness can change everything.",
	"Don’t let yesterday take up too much of today; focus on now.",
	"Success is a journey, not a destination; enjoy every step.",
	"Your mindset determines your success; think positively and take action.",
	"Life is too short to waste on negativity; focus on what matters.",
	"Be fearless in the pursuit of what sets your soul on fire.",
	"Every day is a chance to rewrite your story; make it count.",
	"Don’t let fear stop you from becoming the best version of yourself.",
	"Your vibe attracts your tribe; stay true to who you are.",
	"Dream big, work hard, stay focused, and never give up.",
	"Small steps every day lead to big results over time; keep going.",
	"Be the reason someone smiles today; kindness is always in style.",
	"Your potential is limitless; don’t let doubt hold you back.",
	"Life is a gift; unwrap it with gratitude and joy.",
	"Stay true to yourself; the right people will love the real you.",
	"Every day is a new beginning; make it count and shine bright.",
	"Don’t wait for the perfect moment; take the moment and make it perfect.",
	"Your attitude determines your direction; choose positivity and watch your life change.",
	"Be so good they can’t ignore you; let your work speak for itself.",
	"Every day is a new opportunity to grow, learn, and become better.",
	"Don’t compare your chapter 1 to someone else’s chapter 20; trust your journey.",
	"Stay true to yourself; the right people will love the real you.",
	"Hard work beats talent when talent doesn’t work hard; keep pushing forward.",
	"Life is better when you’re laughing; find joy in the little things.",
	"Believe in yourself, even when no one else does; you’ve got this.",
	"Success is the sum of small efforts repeated day in and day out.",
	"Don’t let fear of failure stop you from chasing what sets your soul on fire.",
	"Be the kind of person who makes others feel seen, heard, and valued.",
	"Your time is limited; don’t waste it living someone else’s life.",
	"Every storm runs out of rain; keep going, brighter days are ahead.",
	"Surround yourself with those who lift you higher and inspire you to grow.",
	"Don’t just exist; live boldly, love deeply, and leave a mark.",
	"Small progress is still progress; celebrate every step forward.",
	"Your potential is endless; don’t let self-doubt hold you back.",
	"Be the energy you want to attract; positivity is contagious.",
	"Life is a canvas; make sure you paint it with vibrant colors.",
	"Don’t wait for opportunities; create them with hard work and determination.",
	"Stay humble, work hard, and be kind; success will follow.",
	"Your dreams are valid; keep working toward them no matter what.",
	"Every setback is a setup for a comeback; keep pushing forward.",
	"Be the light in someone’s darkness; kindness can change everything.",
	"Don’t let yesterday take up too much of today; focus on now.",
	"Success is a journey, not a destination; enjoy every step.",
	"Your mindset determines your success; think positively and take action.",
	"Life is too short to waste on negativity; focus on what matters.",
	"Be fearless in the pursuit of what sets your soul on fire.",
	"Every day is a chance to rewrite your story; make it count.",
	"Don’t let fear stop you from becoming the best version of yourself.",
	"Your vibe attracts your tribe; stay true to who you are.",
	"Dream big, work hard, stay focused, and never give up.",
	"Small steps every day lead to big results over time; keep going.",
	"Be the reason someone smiles today; kindness is always in style.",
	"Your potential is limitless; don’t let doubt hold you back.",
	"Life is a gift; unwrap it with gratitude and joy.",
	"Stay true to yourself; the right people will love the real you.",
	"Every day is a new beginning; make it count and shine bright.",
	"Don’t wait for the perfect moment; take the moment and make it perfect.",
	"Your attitude determines your direction; choose positivity and watch your life change.",
	"Be so good they can’t ignore you; let your work speak for itself.",
	"Every day is a new opportunity to grow, learn, and become better.",
	"Don’t compare your chapter 1 to someone else’s chapter 20; trust your journey.",
	"Stay true to yourself; the right people will love the real you.",
	"Hard work beats talent when talent doesn’t work hard; keep pushing forward.",
	"Life is better when you’re laughing; find joy in the little things.",
	"Believe in yourself, even when no one else does; you’ve got this.",
	"Success is the sum of small efforts repeated day in and day out.",
	"Don’t let fear of failure stop you from chasing what sets your soul on fire.",
	"Be the kind of person who makes others feel seen, heard, and valued.",
	"Your time is limited; don’t waste it living someone else’s life.",
	"Every storm runs out of rain; keep going, brighter days are ahead.",
	"Surround yourself with those who lift you higher and inspire you to grow.",
	"Don’t just exist; live boldly, love deeply, and leave a mark.",
	"Small progress is still progress; celebrate every step forward.",
	"Your potential is endless; don’t let self-doubt hold you back.",
	"Be the energy you want to attract; positivity is contagious.",
	"Life is a canvas; make sure you paint it with vibrant colors.",
	"Don’t wait for opportunities; create them with hard work and determination.",
	"Stay humble, work hard, and be kind; success will follow.",
	"Your dreams are valid; keep working toward them no matter what.",
	"Every setback is a setup for a comeback; keep pushing forward.",
	"Be the light in someone’s darkness; kindness can change everything.",
	"Don’t let yesterday take up too much of today; focus on now.",
	"Success is a journey, not a destination; enjoy every step.",
	"Your mindset determines your success; think positively and take action.",
	"Life is too short to waste on negativity; focus on what matters.",
	"Be fearless in the pursuit of what sets your soul on fire.",
	"Every day is a chance to rewrite your story; make it count.",
	"Don’t let fear stop you from becoming the best version of yourself.",
	"Your vibe attracts your tribe; stay true to who you are.",
	"Dream big, work hard, stay focused, and never give up.",
	"Small steps every day lead to big results over time; keep going.",
	"Be the reason someone smiles today; kindness is always in style.",
	"Your potential is limitless; don’t let doubt hold you back.",
	"Life is a gift; unwrap it with gratitude and joy.",
	"Stay true to yourself; the right people will love the real you.",
	"Every day is a new beginning; make it count and shine bright.",
	"Don’t wait for the perfect moment; take the moment and make it perfect.",
	"Your attitude determines your direction; choose positivity and watch your life change.",
	"Be so good they can’t ignore you; let your work speak for itself.",
	"Every day is a new opportunity to grow, learn, and become better.",
	"Don’t compare your chapter 1 to someone else’s chapter 20; trust your journey.",
	"Stay true to yourself; the right people will love the real you.",
	"Hard work beats talent when talent doesn’t work hard; keep pushing forward.",
	"Life is better when you’re laughing; find joy in the little things.",
	"Believe in yourself, even when no one else does; you’ve got this.",
	"Success is the sum of small efforts repeated day in and day out.",
	"Don’t let fear of failure stop you from chasing what sets your soul on fire.",
	"Be the kind of person who makes others feel seen, heard, and valued.",
	"Your time is limited; don’t waste it living someone else’s life.",
	"Every storm runs out of rain; keep going, brighter days are ahead.",
	"Surround yourself with those who lift you higher and inspire you to grow.",
	"Don’t just exist; live boldly, love deeply, and leave a mark.",
	"Small progress is still progress; celebrate every step forward.",
	"Your potential is endless; don’t let self-doubt hold you back.",
	"Be the energy you want to attract; positivity is contagious.",
	"Life is a canvas; make sure you paint it with vibrant colors.",
	"Don’t wait for opportunities; create them with hard work and determination.",
	"Stay humble, work hard, and be kind; success will follow.",
	"Your dreams are valid; keep working toward them no matter what.",
	"Every setback is a setup for a comeback; keep pushing forward.",
	"Be the light in someone’s darkness; kindness can change everything.",
	"Don’t let yesterday take up too much of today; focus on now.",
	"Success is a journey, not a destination; enjoy every step.",
	"Your mindset determines your success; think positively and take action.",
	"Life is too short to waste on negativity; focus on what matters.",
	"Be fearless in the pursuit of what sets your soul on fire.",
	"Every day is a chance to rewrite your story; make it count.",
	"Don’t let fear stop you from becoming the best version of yourself.",
	"Your vibe attracts your tribe; stay true to who you are.",
	"Dream big, work hard, stay focused, and never give up.",
	"Small steps every day lead to big results over time; keep going.",
	"Be the reason someone smiles today; kindness is always in style.",
	"Your potential is limitless; don’t let doubt hold you back.",
	"Life is a gift; unwrap it with gratitude and joy.",
	"Stay true to yourself; the right people will love the real you.",
	"Every day is a new beginning; make it count and shine bright.",
	"Don’t wait for the perfect moment; take the moment and make it perfect.",
	"Your attitude determines your direction; choose positivity and watch your life change.",
	"Be so good they can’t ignore you; let your work speak for itself.",
	"Every day is a new opportunity to grow, learn, and become better.",
	"Don’t compare your chapter 1 to someone else’s chapter 20; trust your journey.",
	"Stay true to yourself; the right people will love the real you.",
	"Hard work beats talent when talent doesn’t work hard; keep pushing forward.",
	"Life is better when you’re laughing; find joy in the little things.",
	"Believe in yourself, even when no one else does; you’ve got this.",
	"Success is the sum of small efforts repeated day in and day out.",
	"Don’t let fear of failure stop you from chasing what sets your soul on fire.",
}

var potentialTitles = []string{
	"10 Tips to Boost Your Productivity Today",
	"How to Stay Motivated When Life Gets Tough",
	"The Power of Positive Thinking: Transform Your Life",
	"5 Habits of Highly Successful People",
	"Simple Ways to Practice Gratitude Daily",
	"Unlock Your Potential: Believe in Yourself",
	"Morning Routines That Set You Up for Success",
	"Overcoming Fear: Steps to Take the Leap",
	"Finding Balance in a Busy World",
	"The Art of Letting Go: Free Yourself",
	"Mindfulness Made Easy: Start Today",
	"Turning Failures into Stepping Stones",
	"Building Confidence: A Step-by-Step Guide",
	"Healthy Habits for a Happier You",
	"Chasing Dreams: How to Start Today",
	"The Importance of Self-Care in a Busy Life",
	"Creating a Life You Love: Tips and Tricks",
	"From Doubt to Determination: Your Journey",
	"Small Changes, Big Results: Transform Your Life",
	"Finding Joy in the Little Things",
	"How to Build Stronger Relationships",
	"Mastering Time Management: Tips for Success",
	"Turning Goals into Reality: A Practical Guide",
	"The Secret to Staying Consistent",
	"Embracing Change: How to Adapt and Thrive",
	"Daily Affirmations to Boost Your Confidence",
	"Overcoming Procrastination: Take Action Now",
	"The Power of a Growth Mindset",
	"Finding Your Passion: Where to Start",
	"How to Stay Focused in a Distracted World",
	"Building Resilience: Bounce Back Stronger",
	"Creating a Vision for Your Future",
	"The Benefits of Journaling Daily",
	"How to Cultivate a Positive Mindset",
	"Turning Setbacks into Comebacks",
	"Living with Purpose: Find Your Why",
	"Simple Ways to Reduce Stress Daily",
	"The Importance of Setting Boundaries",
	"How to Stay Inspired Every Day",
	"Building a Life of Gratitude and Joy",
	"From Dreaming to Doing: Take the First Step",
	"The Power of Small Wins",
	"How to Stay True to Yourself",
	"Finding Peace in a Chaotic World",
	"Unlocking Creativity: Tips and Techniques",
	"The Art of Saying No Gracefully",
	"Building a Supportive Network",
	"How to Stay Grounded in Tough Times",
	"Turning Challenges into Opportunities",
	"The Joy of Giving: Why It Matters",
	"How to Stay Curious and Keep Learning",
	"Creating a Positive Environment at Home",
	"The Importance of Celebrating Small Wins",
	"How to Stay Patient in a Fast-Paced World",
	"Finding Your Inner Strength",
	"The Power of a Morning Walk",
	"How to Stay Organized and Productive",
	"Building a Life of Meaning and Fulfillment",
	"The Benefits of a Digital Detox",
	"How to Stay Hopeful in Difficult Times",
	"Turning Ideas into Action: A Step-by-Step Guide",
	"The Importance of Self-Reflection",
	"How to Stay Energized Throughout the Day",
	"Building a Life You’re Proud Of",
	"The Power of a Smile: Spread Positivity",
	"How to Stay Focused on Your Goals",
	"Finding Calm in the Chaos",
	"The Joy of Simple Living",
	"How to Stay Grateful Every Day",
	"Building a Stronger Mindset",
	"The Importance of Taking Breaks",
	"How to Stay Inspired by Others",
	"Turning Obstacles into Opportunities",
	"The Power of a Positive Attitude",
	"How to Stay Committed to Your Goals",
	"Finding Happiness in the Present Moment",
	"The Benefits of a Daily Routine",
	"How to Stay Motivated Long-Term",
	"Building a Life of Balance and Harmony",
	"The Power of Visualization: See Your Success",
	"How to Stay True to Your Values",
	"Finding Your Flow: Tips for Peak Performance",
	"The Importance of Lifelong Learning",
	"How to Stay Resilient in Tough Times",
	"Turning Dreams into Reality: A Practical Guide",
	"The Power of a Supportive Community",
	"How to Stay Focused on What Matters",
	"Finding Joy in Everyday Moments",
	"The Benefits of a Positive Mindset",
	"How to Stay Grounded and Centered",
	"Building a Life of Purpose and Passion",
	"The Power of Small, Consistent Actions",
	"How to Stay Inspired by Your Own Journey",
	"Turning Fear into Fuel for Success",
	"The Importance of Self-Discipline",
	"How to Stay Grateful in Challenging Times",
	"Building a Life of Gratitude and Abundance",
	"The Power of a Clear Vision",
	"How to Stay Motivated When Progress is Slow",
	"Finding Peace in the Present Moment",
	"The Joy of Living with Intention",
	"How to Stay Focused on Your Priorities",
	"Building a Life of Joy and Fulfillment",
	"The Power of a Positive Morning Routine",
	"How to Stay True to Your Dreams",
	"Finding Strength in Vulnerability",
	"The Importance of Celebrating Your Wins",
	"How to Stay Inspired by the World Around You",
	"Turning Challenges into Growth Opportunities",
	"The Power of a Grateful Heart",
	"How to Stay Committed to Your Vision",
	"Building a Life of Love and Connection",
	"The Joy of Living Authentically",
	"How to Stay Focused on Your Purpose",
	"Finding Peace in Letting Go",
	"The Power of a Supportive Mindset",
	"How to Stay Motivated by Your Why",
	"Building a Life of Gratitude and Joy",
	"The Importance of Taking Small Steps",
	"How to Stay Inspired by Your Own Progress",
	"Turning Setbacks into Stepping Stones",
	"The Power of a Positive Outlook",
	"How to Stay True to Your Path",
	"Finding Joy in the Journey",
	"The Benefits of Living with Gratitude",
	"How to Stay Focused on Your Dreams",
	"Building a Life of Meaning and Impact",
	"The Power of a Clear and Focused Mind",
	"How to Stay Inspired by Your Own Growth",
	"Turning Obstacles into Opportunities for Success",
	"The Importance of Staying True to Yourself",
	"How to Stay Grateful in Every Moment",
	"Building a Life of Purpose and Passion",
	"The Joy of Living with Intention and Gratitude",
}

var potentialTags = []string{
	"#Motivation", "#Inspiration", "#SelfCare", "#Mindfulness", "#Growth", "#Success", "#Positivity", "#Gratitude", "#Happiness", "#Goals", "#Focus", "#Resilience", "#Balance", "#DreamBig", "#HealthyHabits", "#Mindset", "#DailyJoy", "#SelfLove", "#Productivity", "#KeepGoing",
}

var potentialComments = []string{
	"Love this!", "So inspiring!", "Needed this today!", "Great advice!", "So true!", "Thank you for sharing!", "This is amazing!", "Love your perspective!", "Spot on!", "You’re awesome!", "So motivating!", "This made my day!", "Perfect timing!", "Love the energy!", "So relatable!", "Keep it up!", "This is gold!", "You nailed it!", "So powerful!", "Love the positivity!", "This is fire!", "So well said!", "You’re crushing it!", "This hit home!", "So uplifting!", "Love your vibe!", "This is everything!", "So encouraging!", "You’re a star!", "This is brilliant!", "So refreshing!", "Love your work!", "This is so true!", "You’re inspiring!", "So motivating!", "This is pure gold!", "Love your mindset!", "This is spot on!", "You’re amazing!", "So empowering!", "This is fantastic!", "Love your energy!", "This is so helpful!", "You’re a legend!", "So inspiring!", "This is incredible!", "Love your passion!", "This is so powerful!", "You’re killing it!", "So uplifting!", "This is perfect!", "Love your creativity!", "So well-written!", "This is a gem!", "You’re on fire!", "So thoughtful!", "This is life-changing!", "Love your authenticity!", "This is so relatable!", "You’re a rockstar!", "So motivating!", "This is pure inspiration!", "Love your dedication!", "This is so true!", "You’re a genius!", "So uplifting!", "This is amazing work!", "Love your insights!", "This is so encouraging!", "You’re a true inspiration!", "So powerful!", "This is exactly what I needed!", "Love your positivity!", "This is so well said!", "You’re incredible!", "So motivating!", "This is gold!", "Love your energy!", "This is so inspiring!", "You’re a blessing!", "So uplifting!", "This is fantastic advice!", "Love your perspective!", "This is so true!", "You’re amazing!", "So empowering!", "This is brilliant!", "Love your vibe!", "This is so helpful!", "You’re a legend!", "So inspiring!", "This is incredible!", "Love your passion!", "This is so powerful!", "You’re killing it!", "So uplifting!", "This is perfect!", "Love your creativity!", "So well-written!", "This is a gem!", "You’re on fire!", "So thoughtful!", "This is life-changing!", "Love your authenticity!", "This is so relatable!", "You’re a rockstar!", "So motivating!", "This is pure inspiration!", "Love your dedication!", "This is so true!", "You’re a genius!", "So uplifting!", "This is amazing work!", "Love your insights!", "This is so encouraging!", "You’re a true inspiration!", "So powerful!", "This is exactly what I needed!", "Love your positivity!", "This is so well said!", "You’re incredible!", "So motivating!", "This is gold!", "Love your energy!", "This is so inspiring!", "You’re a blessing!", "So uplifting!", "This is fantastic advice!", "Love your perspective!", "This is so true!", "You’re amazing!", "So empowering!", "This is brilliant!", "Love your vibe!", "This is so helpful!", "You’re a legend!", "So inspiring!", "This is incredible!", "Love your passion!", "This is so powerful!", "You’re killing it!", "So uplifting!", "This is perfect!", "Love your creativity!", "So well-written!", "This is a gem!", "You’re on fire!", "So thoughtful!", "This is life-changing!", "Love your authenticity!", "This is so relatable!", "You’re a rockstar!", "So motivating!", "This is pure inspiration!", "Love your dedication!", "This is so true!", "You’re a genius!", "So uplifting!", "This is amazing work!", "Love your insights!", "This is so encouraging!", "You’re a true inspiration!", "So powerful!", "This is exactly what I needed!", "Love your positivity!", "This is so well said!", "You’re incredible!", "So motivating!", "This is gold!", "Love your energy!",
}

func Seed(store store.Storage, numUsers int, numPosts int, numComments int, db *sql.DB)  {

	// seed users
	users := generateUsers(numUsers)

	ctx := context.Background()

	tx, _ := db.BeginTx(ctx, nil)
	
	for i := range numUsers {
		if err := store.Users.Create(ctx, tx, users[i]); err != nil {
			log.Fatal(err)
		}

	}

	tx.Commit()
	// seed posts
	posts := generatePosts(numPosts, users)

	for i := range numPosts {
		if err := store.Posts.Create(ctx, posts[i]); err != nil {
			log.Fatal(err)
		}

	}

	// seed comments
	comments := generateComments(numComments, users, posts)
	for i := range numComments {

		if err := store.Comments.Create(ctx, comments[i]); err != nil {
			log.Fatal(err)
		}

	}
	log.Println("Seeding Complete !")
}

func generateUsers(numUsers int) []*store.User {

	users := make([]*store.User, numUsers)

	for i := 0; i < numUsers; i++ {
		name := potentialNames[rand.Intn(len(potentialNames))]+strconv.Itoa(i)

		users[i] = &store.User{
			Username: name,
			Email:    fmt.Sprintf("%s@example.com", name),
			
		}
		_ = users[i].Password.Set(name)
	}

	return users

}

func generatePosts(numPosts int, users []*store.User) []*store.Post {

	posts := make([]*store.Post, numPosts)

	for i := range numPosts {
		posts[i] = &store.Post{
			Title:   potentialTitles[rand.Intn(len(potentialTitles))],
			Content: potentialContents[rand.Intn(len(potentialContents))],
			UserID:  users[rand.Intn(len(users))].ID,
			Tags: []string{
				potentialTags[rand.Intn(len(potentialTags))],
				potentialTags[rand.Intn(len(potentialTags))],
			},
		}

	}

	return posts
}

func generateComments(numComments int, users []*store.User, posts []*store.Post) []*store.Comment {

	comments := make([]*store.Comment, numComments)

	for i := range numComments {
		comments[i] = &store.Comment{
			PostID:  posts[rand.Intn(len(posts))].ID,
			UserID:  users[rand.Intn(len(users))].ID,
			Content: potentialComments[rand.Intn(len(potentialComments))],
		}
	}

	return comments
}
