#!/usr/bin/ruby
#file: main.rb

require_relative 'search'
require 'mail'
require 'net/smtp'
require 'mysql2'
require 'dotenv'

Dotenv.load

search	= Search.new('locations', 'alsace')

search.min_price 	= 450
search.max_price 	= 600
search.min_area 	= 40
search.min_room	 	= 2
search.type 		= 2
search.furniture 	= 2

search.locations = [
	'Bischheim', 
	'Hoenheim', 
	'Souffelweyersheim',
	'Reichstett',
	'Schiltigheim',
	'Mundolsheim',
	'Lampertheim',
	'Vendenheim',
	'Berstett',
	'Weyersheim',
	'Brumath',
	'Niederhausbergen',
	'Oberhausbergen',
	'Mittelhausbergen',
	'Strasbourg',
	'Ostwald',
	'Gambsheim',
	'Hoerdt',
	'Kilstett',
	'La%20Wantzenau',
	'Eckbolsheim'
]

results = search.results

if results.count > 0 then

	begin
		client = Mysql2::Client.new(:host => ENV['DB_HOST'], :username => ENV['DB_USERNAME'], :password => ENV['DB_PASSWORD'], :database => ENV['DB_DATABASE'])

		rows = client.query("SELECT * FROM #{ENV['DB_TABLE']}")

		rows.each do |row|
			if results.include? row['url']
				results.delete(row['url'])
			end
		end

		if results.count > 0 then
			query = "INSERT INTO annonces (url, created_at) VALUES "

			results.each_with_index do |result, index|
				query.concat("(\'#{result}\', NOW())")

				if index.to_i < results.count-1 then
					query.concat(", ")
				else
					query.concat(";")
				end
			end

			client.query(query)
		end
		
	rescue Exception => e
		puts e.inspect

		puts ENV['DB_HOST']
		puts ENV['SMTP_DOMAIN']
		puts ENV['MAIL_FROM']
	ensure
		client.close if client
	end
	
	if results.count > 0 then
		message = search.criteria
		message.concat("<br/><b>Results</b><br/><ul>")
		
		results.each do |result|
			message.concat("<li><a href='#{result}'>#{result}</a></li>")
		end

		message.concat("</ul>")

		Mail.defaults do
			delivery_method :smtp, {
				:address		=> ENV['SMTP_HOST'],
				:port			=> ENV['SMTP_PORT'],
				:domain			=> ENV['SMTP_DOMAIN'],
				:user_name		=> ENV['SMTP_USERNAME'],
				:password		=> ENV['SMTP_PASSWORD'],
				:authentication		=> 'plain',
				:enable_starttls_auto	=> true
			}
		end

		Mail.deliver do
			to	ENV['MAIL_TO']
			from	ENV['MAIL_FROM']
			subject "New Appartments (#{results.count})"
		
			html_part do
				content_type 'text/html; charset_UTF-8'
				body message
			end
		end		
	end
end
