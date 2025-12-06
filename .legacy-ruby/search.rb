#!/usr/bin/ruby
#file: search.rb

require 'nokogiri'
require 'open-uri'

class Search
	
	attr_accessor :min_price, :max_price, :min_area, :max_area, :min_room, :max_room, :type, :furniture

	@@url
	@@locations

	def initialize(category, region)
		@url = "http://www.leboncoin.fr/#{category}/offres/#{region}/?"
	end

	def locations=(locations)
		@locations = locations
	end

	def criteria
		criteria = "<b>Criteria</b><br/><ul>"

		if @min_price then
			criteria.concat("<li>Min Price : #{@min_price} euros</li>")
                end

                if @max_price then
                       	criteria.concat("<li>Max Price : #{@max_price} euros</li>")
                end

                if @min_area then
                        criteria.concat("<li>Min Area : #{@min_area} m2</li>")
                end

                if @max_area then
                        criteria.concat("<li>Max Area : #{@max_area} m2</li>")
                end

                if @min_room then
                      	criteria.concat("<li>Min Number of Room : #{@min_room}</li>")
                end

                if @max_room then
                        criteria.concat("<li>Max Number of Room : #{@max_room}</li>")
                end

                if @type then
                        criteria.concat("<li>Type : #{@type}</li>")
                end

                if @furniture then
                       	criteria.concat("<li>Furniture : #{@furniture}</li>")
                end

		if @locations then
			criteria.concat("<li>Locations : #{@locations.join(', ')}</li>")
		end

		criteria.concat("</ul>")

		criteria
	end

	def url(location)
		tmp = @url + "location=#{location}"

		if @min_price then
			tmp.concat("&mrs=#{@min_price}")
		end

		if @max_price then
			tmp.concat("&mre=#{@max_price}")
		end

		if @min_area then
			tmp.concat("&sqs=#{@min_area}")
		end

		if @max_area then
			tmp.concat("&sqe=#{@max_area}")
		end

		if @min_room then
			tmp.concat("&ros=#{@min_room}")
		end

		if @max_room then
			tmp.concat("&roe=#{@max_room}")
		end

		if @type then
			tmp.concat("&ret=#{@type}")
		end

		if @furniture then
			tmp.concat("&furn=#{@furniture}")
		end

		tmp		
	end

	def results
		results = Array.new
		if @locations then
			@locations.each do |location|
			
				url = self.url(location)
				html = Nokogiri::HTML(open(url))

				html.css('div.list-lbc a').map do |link| 
					link.css('div.lbc div.date div:contains("Aujourd")').map do |date|
						if link['href'].match("/locations/") then
							results.push(link['href'])
						end
					end
				end
			end
		end

		results.uniq

		results
	end
end
