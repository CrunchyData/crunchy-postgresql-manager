define([
  'angular',
  'lodash',
  'kbn',
  'moment',
  './queryCtrl',
],
function (angular, _, kbn) {
  'use strict';

  var module = angular.module('grafana.services');

  module.factory('PrometheusDatasource', function($q, $http, templateSrv) {

    function PrometheusDatasource(datasource) {
      this.type = 'prometheus';
      this.editorSrc = 'app/features/prometheus/partials/query.editor.html';
      this.name = datasource.name;
      this.supportMetrics = true;

      var url = datasource.url;
      if (url[url.length-1] === '/') {
        // remove trailing slash
        url = url.substr(0, url.length - 1);
      }
      this.url = url;
      this.lastErrors = {};
    }

    // Called once per panel (graph)
    PrometheusDatasource.prototype.query = function(options) {
      var end = convertToPrometheusTime(options.range.to);
      var range = convertToPrometheusRange(options.range.from, options.range.to);

      var queries = [];
      _.each(options.targets, _.bind(function(target) {
        if (!target.expr || target.hide) {
          return;
        }

        var query = {};
        query.expr = templateSrv.replace(target.expr, options.scopedVars);

        var interval = target.interval || options.interval;
        var intervalFactor = target.intervalFactor || 1;
        query.step = this.calculateInterval(interval, intervalFactor);

        queries.push(query);
      }, this));

      // No valid targets, return the empty result to save a round trip.
      if (_.isEmpty(queries)) {
        var d = $q.defer();
        d.resolve({ data: [] });
        return d.promise;
      }

      var allQueryPromise = _.map(queries, _.bind(function(query) {
        return this.performTimeSeriesQuery(query, range, end);
      }, this));

      var self = this;
      return $q.all(allQueryPromise)
        .then(function(allResponse) {
          var result = [];

          _.each(allResponse, function(response, index) {
            if (response.data.type === 'error') {
              self.lastErrors.query = response.data.value;
              throw response.data.value;
            }
            delete self.lastErrors.query;

            _.each(response.data.value, function(metricData) {
              result.push(transformMetricData(metricData, options.targets[index]));
            });
          });

          return { data: result };
        });
    };

    PrometheusDatasource.prototype.performTimeSeriesQuery = function(query, range, end) {
      var url = this.url + '/api/query_range?expr=' + encodeURIComponent(query.expr) + '&range=' + range + '&end=' + end;

      var step = query.step;
      // Prometheus drop query if range/step > 11000
      // calibrate step if it is too big
      if (step !== 0 && range / step > 11000) {
        step = Math.floor(range / 11000);
      }
      url += '&step=' + step;

      var options = {
        method: 'GET',
        url: url,
      };

      return $http(options);
    };

    PrometheusDatasource.prototype.performSuggestQuery = function(query) {
      var options = {
        method: 'GET',
        url: this.url + '/api/metrics',
      };

      return $http(options).then(function(result) {
        var suggestData = _.filter(result.data, function(metricName) {
          return metricName.indexOf(query) !==  1;
        });

        return suggestData;
      });
    };

    PrometheusDatasource.prototype.metricFindQuery = function(query) {
      var options;
      var matches = query.match(/^[a-zA-Z_:*][a-zA-Z0-9_:*]*/);

      if (matches != null && matches[0].indexOf('*') >= 0) {
        // if query has wildcard character, return metric name list
        options = {
          method: 'GET',
          url: this.url + '/api/metrics',
        };

        return $http(options)
          .then(function(result) {
            return _.chain(result.data)
              .filter(function(metricName) {
                var r = new RegExp(matches[0].replace(/\*/g, '.*'));
                return r.test(metricName);
              })
              .map(function(matchedMetricName) {
                return {
                  text: matchedMetricName,
                  expandable: true
                };
              })
              .value();
            });
      } else {
        // if query contains full metric name, return metric name and label list
        options = {
          method: 'GET',
          url: this.url + '/api/query?expr=' + encodeURIComponent(query),
        };

        return $http(options)
          .then(function(result) {
            return _.map(result.data.value, function(metricData) {
              return {
                text: getOriginalMetricName(metricData.metric),
                expandable: true
              };
            });
          });
        }
    };

    PrometheusDatasource.prototype.calculateInterval = function(interval, intervalFactor) {
      var sec = kbn.interval_to_seconds(interval);

      if (sec < 1) {
        sec = 1;
      }

      return sec * intervalFactor;
    };

    function transformMetricData(md, options) {
      var dps = [],
          metricLabel = null;

      metricLabel = createMetricLabel(md.metric, options);

      dps = _.map(md.values, function(value) {
        return [parseFloat(value[1]), value[0] * 1000];
      });

      return { target: metricLabel, datapoints: dps };
    }

    function createMetricLabel(labelData, options) {
      if (_.isUndefined(options) || _.isEmpty(options.legendFormat)) {
        return getOriginalMetricName(labelData);
      }

      var originalSettings = _.templateSettings;
      _.templateSettings = {
        interpolate: /\{\{(.+?)\}\}/g
      };

      var template = _.template(options.legendFormat);
      var metricName = template(labelData);

      _.templateSettings = originalSettings;

      return metricName;
    }

    function getOriginalMetricName(labelData) {
      var metricName = labelData.__name__ || '';
      delete labelData.__name__;
      var labelPart = _.map(_.pairs(labelData), function(label) {
        return label[0] + '="' + label[1] + '"';
      }).join(',');
      return metricName + '{' + labelPart + '}';
    }

    function convertToPrometheusRange(from, to) {
      return Math.floor(convertToPrometheusTime(to) - convertToPrometheusTime(from));
    }

    function convertToPrometheusTime(date) {
      date = kbn.parseDate(date);
      return date.getTime() / 1000;
    }

    return PrometheusDatasource;
  });

});
