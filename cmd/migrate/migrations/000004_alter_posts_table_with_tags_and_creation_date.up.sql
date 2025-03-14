ALTER TABLE posts
ADD 
    COLUMN tags varchar(100)[];

ALTER TABLE posts
ADD 
    COLUMN     Update_date timestamp(0) with time zone NOT NULL DEFAULT NOW();

