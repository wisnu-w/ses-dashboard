import { useState, useEffect, useCallback, useRef } from 'react';
import { Search, Calendar, Filter, X } from 'lucide-react';
import Layout from '../components/Layout';
import EventsTable from '../components/EventsTable';
import { eventsService } from '../services/api';
import type { EventsResponse } from '../types/api';

// Custom hook for debouncing
const useDebounce = (value: string, delay: number) => {
  const [debouncedValue, setDebouncedValue] = useState(value);

  useEffect(() => {
    const handler = setTimeout(() => {
      setDebouncedValue(value);
    }, delay);

    return () => {
      clearTimeout(handler);
    };
  }, [value, delay]);

  return debouncedValue;
};

const EventsPage = () => {
  const [eventsData, setEventsData] = useState<EventsResponse | null>(null);
  const [loading, setLoading] = useState(true);
  const [currentPage, setCurrentPage] = useState(1);
  const [pageSize] = useState(25);
  const [search, setSearch] = useState('');
  const [startDate, setStartDate] = useState('');
  const [endDate, setEndDate] = useState('');
  const [showFilters, setShowFilters] = useState(false);
  
  // Debounce search input to reduce API calls
  const debouncedSearch = useDebounce(search, 800);
  const lastRequestKeyRef = useRef('');
  const requestSequenceRef = useRef(0);

  const loadEvents = useCallback(async (page = 1, searchTerm = debouncedSearch, start = startDate, end = endDate) => {
    const requestKey = `${page}|${pageSize}|${searchTerm}|${start}|${end}`;
    if (lastRequestKeyRef.current === requestKey) return;

    lastRequestKeyRef.current = requestKey;
    const requestSequence = requestSequenceRef.current + 1;
    requestSequenceRef.current = requestSequence;

    try {
      setLoading(true);
      const params = new URLSearchParams({
        page: page.toString(),
        limit: pageSize.toString()
      });
      
      if (searchTerm) params.append('search', searchTerm);
      if (start) params.append('start_date', start);
      if (end) params.append('end_date', end);
      
      const response = await fetch(`/api/events?${params}`, {
        headers: {
          'Authorization': `Bearer ${localStorage.getItem('token')}`
        }
      });
      
      if (!response.ok) {
        throw new Error(`HTTP error! status: ${response.status}`);
      }
      
      const data = await response.json();
      if (requestSequenceRef.current !== requestSequence) return;
      setEventsData(data);
      setCurrentPage(page);
    } catch (error) {
      console.error('Failed to load events:', error);
      // Fallback to eventsService if direct fetch fails
      try {
        const data = await eventsService.getEvents(page, pageSize, searchTerm, start, end);
        if (requestSequenceRef.current !== requestSequence) return;
        setEventsData(data);
        setCurrentPage(page);
      } catch (fallbackError) {
        console.error('Fallback also failed:', fallbackError);
      }
    } finally {
      if (requestSequenceRef.current === requestSequence) {
        setLoading(false);
      }
    }
  }, [debouncedSearch, startDate, endDate, pageSize]);

  const handleSearch = () => {
    setCurrentPage(1);
    loadEvents(1, search, startDate, endDate);
  };

  const clearFilters = () => {
    setSearch('');
    setStartDate('');
    setEndDate('');
    setCurrentPage(1);
    lastRequestKeyRef.current = '';
  };

  useEffect(() => {
    setCurrentPage(1);
    loadEvents(1, debouncedSearch, startDate, endDate);
  }, [debouncedSearch, startDate, endDate, loadEvents]);

  const handlePageChange = (page: number) => {
    setCurrentPage(page);
    loadEvents(page, debouncedSearch, startDate, endDate);
  };

  return (
    <Layout title="Events">
      <div className="space-y-6">
        <div className="flex items-center justify-between">
          <div>
            <h1 className="text-2xl font-bold text-gray-900">Email Events</h1>
            <p className="text-gray-600 mt-1">Detailed view of all SES email events</p>
          </div>
          <button
            onClick={() => setShowFilters(!showFilters)}
            className="inline-flex items-center px-4 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700 transition-colors duration-200"
          >
            <Filter className="w-4 h-4 mr-2" />
            {showFilters ? 'Hide Filters' : 'Show Filters'}
          </button>
        </div>

        {/* Search and Filter Section */}
        {showFilters && (
          <div className="bg-white rounded-lg shadow-sm border border-gray-200 p-6">
            <div className="grid grid-cols-1 md:grid-cols-4 gap-4">
              <div className="md:col-span-2">
                <label className="block text-sm font-medium text-gray-700 mb-2">Search</label>
                <div className="relative">
                  <Search className="absolute left-3 top-1/2 transform -translate-y-1/2 text-gray-400 w-4 h-4" />
                  <input
                    type="text"
                    value={search}
                    onChange={(e) => setSearch(e.target.value)}
                    placeholder="Search by email, subject, or source... (auto-search)"
                    className="w-full pl-10 pr-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
                  />
                </div>
              </div>
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-2">Start Date</label>
                <div className="relative">
                  <Calendar className="absolute left-3 top-1/2 transform -translate-y-1/2 text-gray-400 w-4 h-4" />
                  <input
                    type="date"
                    value={startDate}
                    onChange={(e) => setStartDate(e.target.value)}
                    className="w-full pl-10 pr-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
                  />
                </div>
              </div>
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-2">End Date</label>
                <div className="relative">
                  <Calendar className="absolute left-3 top-1/2 transform -translate-y-1/2 text-gray-400 w-4 h-4" />
                  <input
                    type="date"
                    value={endDate}
                    onChange={(e) => setEndDate(e.target.value)}
                    className="w-full pl-10 pr-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
                  />
                </div>
              </div>
            </div>
            <div className="flex space-x-3 mt-4">
              <button
                onClick={handleSearch}
                className="px-4 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700 transition-colors duration-200 flex items-center"
              >
                <Search className="w-4 h-4 mr-2" />
                Search
              </button>
              <button
                onClick={clearFilters}
                className="px-4 py-2 bg-gray-300 text-gray-700 rounded-lg hover:bg-gray-400 transition-colors duration-200 flex items-center"
              >
                <X className="w-4 h-4 mr-2" />
                Clear
              </button>
            </div>
          </div>
        )}

        {(loading || (eventsData && eventsData.events && eventsData.events.length > 0)) && (
          <EventsTable
            events={eventsData?.events ?? []}
            pagination={eventsData?.pagination ?? {
              page: currentPage,
              limit: pageSize,
              total: 0,
              totalPages: 1,
              hasNext: false,
              hasPrev: false,
            }}
            onPageChange={handlePageChange}
            loading={loading}
          />
        )}

        {eventsData && eventsData.events && eventsData.events.length === 0 && !loading && (
          <div className="bg-white rounded-lg shadow-sm border border-gray-200 p-8">
            <div className="text-center">
              <Search className="w-16 h-16 text-gray-300 mx-auto mb-4" />
              <h3 className="text-lg font-medium text-gray-900 mb-2">No events found</h3>
              <p className="text-gray-600 mb-4">
                {search || startDate || endDate 
                  ? 'No events match your search criteria. Try adjusting your filters.' 
                  : 'No events available at the moment.'}
              </p>
              {(search || startDate || endDate) && (
                <button
                  onClick={clearFilters}
                  className="px-4 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700 transition-colors duration-200"
                >
                  Clear Filters
                </button>
              )}
            </div>
          </div>
        )}
      </div>
    </Layout>
  );
};

export default EventsPage;
